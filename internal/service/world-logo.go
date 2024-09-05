package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"ext-data-domain/internal/model"
	sq "ext-data-domain/internal/service/query"

	"github.com/oklog/ulid/v2"
	"github.com/slyngshot-al/packages/log"
	"github.com/slyngshot-al/packages/storage"
	"github.com/slyngshot-al/packages/storage/psql"
)

type WorldLogoServiceConfig struct {
	ApiKey      string `yaml:"apikey"`
	DbConn      storage.MainConnection
	FileStorage FileStorage
}

type WorldLogoService struct {
	initTx      storage.InitTxFn
	initCtxConn storage.InitCtxConnFn
	apiKey      string
	fileStorage FileStorage
}

func NewWorldLogoService(conf WorldLogoServiceConfig) *WorldLogoService {
	return &WorldLogoService{
		initTx:      storage.InitTxFnFromConn(conf.DbConn),
		initCtxConn: storage.InitCtxConnFnFromConn(conf.DbConn),
		apiKey:      conf.ApiKey,
		fileStorage: conf.FileStorage,
	}
}

//go:generate mockery --name=FileStorage --structname=FileStorage --output=./mocks/ --case=underscore
type FileStorage interface {
	GetBaseUrl() string
	Upload(ctx context.Context, filename string, body io.Reader, contentType string) (string, error)
	CopyFolder(ctx context.Context, sourceFolder, destinationFolder string) error
	Delete(ctx context.Context, key string) error
}

func (s *WorldLogoService) GetWorldLogoById(ctx context.Context, id string) (res model.WorldLogo, err error) {
	mtxFn := CollectMetricFn("GetWorldLogoById")
	defer func() {
		mtxFn(ctx, err)
	}()
	ctx = log.CtxWithValues(ctx, "action", "GetWorldLogoById", "id", id)
	ctx = s.initCtxConn(ctx)
	res, err = sq.GetWorldLogoById(ctx, id)
	res.LogoPath = model.UrlJoinPath(s.fileStorage.GetBaseUrl(), res.LogoPath)
	return res, fromStorageErr(err)
}

func (s *WorldLogoService) SaveWorldLogo(ctx context.Context, apiKey string, input model.SaveWorldLogoInput) (id string, err error) {
	mtxFn := CollectMetricFn("SaveWorldLogo")
	defer func() {
		mtxFn(ctx, err)
	}()

	if apiKey != s.apiKey {
		return "", ErrForbidden
	}

	// validate input
	if strings.TrimSpace(input.SrcKey) == "" {
		return "", ErrInvalidParams.WithMessage("'src_key' must be provided")
	}
	if strings.TrimSpace(input.Name) == "" {
		return "", ErrInvalidParams.WithMessage("'name' must be provided")
	}
	if len(input.LogoBase64Str) == 0 {
		return "", ErrInvalidParams.WithMessage("'logo_base64_str' must be provided")
	}
	if len(input.ContentType) == 0 {
		return "", ErrInvalidParams.WithMessage("'content_type' must be provided")
	}
	if len(input.FileExtension) == 0 {
		return "", ErrInvalidParams.WithMessage("'file_extension' must be provided")
	}

	ctx = log.CtxWithValues(ctx, "action", "SaveWorldLogo", "name", input.Name, "key", input.SrcKey)
	ctx = s.initCtxConn(ctx)

	rec := model.WorldLogoInput{
		Id:     ulid.Make().String(),
		Name:   input.Name,
		SrcKey: input.SrcKey,
	}

	err = storage.DoTransactionAction(ctx, s.initTx, func(ctx context.Context) (err error) {
		// overwrite id if src_key exists
		cRec, err := sq.LockWorldLogoBySrcKey(ctx, input.SrcKey)
		if err != nil {
			if !isStorageNotFoundErr(err) {
				return fromStorageErr(err)
			}
		}
		if cRec.Id != "" {
			rec.Id = cRec.Id
		}

		// upload logo file
		path, err := s.doUploadLogo(ctx, fmt.Sprintf("worldlogo/%s.%s", rec.Id, input.FileExtension), input.ContentType, []byte(input.LogoBase64Str))
		if err != nil {
			return err
		}
		rec.LogoPath = path

		// save record
		if err = sq.SaveWorldLogo(ctx, rec); err != nil {
			return fromStorageErr(err)
		}
		return nil
	})
	return rec.Id, err
}

func (s *WorldLogoService) DeleteWorldLogo(ctx context.Context, apiKey string, ids ...string) (err error) {
	mtxFn := CollectMetricFn("DeleteWorldLogo")
	defer func() {
		mtxFn(ctx, err)
	}()

	if apiKey != s.apiKey {
		return ErrForbidden
	}

	ctx = log.CtxWithValues(ctx, "action", "DeleteWorldLogo", "ids", ids)

	return storage.DoTransactionAction(ctx, s.initTx, func(ctx context.Context) (err error) {
		recs, err := sq.GetWorldLogos(ctx, model.WorldLogosQueryOptions{Ids: ids}, nil, psql.Pagination{Limit: uint64(len(ids))})
		if err != nil {
			return fromStorageErr(err)
		}
		if err = sq.DeleteWorldLogo(ctx, ids...); err != nil {
			return fromStorageErr(err)
		}
		// try to clean up the files
		go func() {
			for _, rec := range recs {
				if rec.LogoPath != "" && !strings.HasPrefix(rec.LogoPath, "http") {
					err = s.fileStorage.Delete(ctx, rec.LogoPath)
					if err != nil {
						log.Error(ctx, err, "failed to delete file", "file", rec.LogoPath)
					}
				}
			}
		}()
		return nil
	})
}

func (s *WorldLogoService) GetWorldLogos(ctx context.Context, ops model.WorldLogosQueryOptions, sort []psql.Sort, pg psql.Pagination) (res []model.WorldLogo, err error) {
	mtxFn := CollectMetricFn("GetWorldLogos")
	defer func() {
		mtxFn(ctx, err)
	}()
	ctx = log.CtxWithValues(ctx, "action", "GetWorldLogos")
	ctx = s.initCtxConn(ctx)
	res, err = sq.GetWorldLogos(ctx, ops, sort, pg)
	for i := range res {
		res[i].LogoPath = model.UrlJoinPath(s.fileStorage.GetBaseUrl(), res[i].LogoPath)
	}
	return res, fromStorageErr(err)
}

func (s *WorldLogoService) GetWorldLogosCount(ctx context.Context, ops model.WorldLogosQueryOptions) (count uint64, err error) {
	mtxFn := CollectMetricFn("GetWorldLogosCount")
	defer func() {
		mtxFn(ctx, err)
	}()
	ctx = log.CtxWithValues(ctx, "action", "GetWorldLogosCount")
	ctx = s.initCtxConn(ctx)
	count, err = sq.GetWorldLogosCount(ctx, ops)
	return count, fromStorageErr(err)
}

func (s *WorldLogoService) doUploadLogo(ctx context.Context, name string, contentType string, data []byte) (path string, err error) {
	source := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))
	propImage, err := s.fileStorage.Upload(ctx, name, source, contentType)
	if err != nil {
		return "", ErrInternal.Consume(err).WithAdditionalInfo("failed to upload file", map[string]any{"logo_base64_str": err.Error()})
	}

	return propImage, nil
}
