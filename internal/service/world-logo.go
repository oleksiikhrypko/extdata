package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"ext-data-domain/internal/model"
	sq "ext-data-domain/internal/service/query"

	"github.com/Slyngshot-Team/packages/log"
	"github.com/Slyngshot-Team/packages/storage"
	"github.com/Slyngshot-Team/packages/storage/psql"
	"github.com/oklog/ulid/v2"
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

func (s *WorldLogoService) SaveWorldLogo(ctx context.Context, apiKey string, input model.WorldLogoInput) (id string, err error) {
	mtxFn := CollectMetricFn("SaveWorldLogo")
	defer func() {
		mtxFn(ctx, err)
	}()

	if apiKey != s.apiKey {
		return "", ErrForbidden
	}

	ctx = log.CtxWithValues(ctx, "action", "SaveWorldLogo", "id", input.Id, "name", input.Name, "key", input.SrcKey)
	ctx = s.initCtxConn(ctx)
	if input.Id == "" {
		input.Id = ulid.Make().String()
	}

	return input.Id, storage.DoTransactionAction(ctx, s.initTx, func(ctx context.Context) (err error) {
		if err = s.doUploadLogo(ctx, &input); err != nil {
			return err
		}

		if err = sq.SaveWorldLogo(ctx, input); err != nil {
			return fromStorageErr(err)
		}
		return nil
	})
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

func (s *WorldLogoService) doUploadLogo(ctx context.Context, item *model.WorldLogoInput) error {
	repl := strings.NewReplacer(" ", "_", "*", "_", "\\", "_", "-", "_", "/", "_")
	name := repl.Replace(item.Name)

	timePfx := time.Now().UnixNano()
	// ValuePropositionImage
	if item.LogoData != nil {
		source := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(*item.LogoData)))
		propImage, err := s.fileStorage.Upload(ctx, fmt.Sprintf("/worldlogo/%s_%d", name, timePfx), source, "")
		if err != nil {
			return ErrInternal.Consume(err).WithAdditionalInfo("failed to upload file", map[string]any{"logo_data": err.Error()})
		}
		item.LogoPath = &propImage
	}

	return nil
}
