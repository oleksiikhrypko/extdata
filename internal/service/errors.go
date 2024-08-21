package service

import (
	"errors"

	"github.com/Slyngshot-Team/packages/storage"
	"github.com/Slyngshot-Team/packages/xerrors"
)

var (
	ErrForbidden          = xerrors.New("forbidden").WithExtensions(map[string]interface{}{"code": "403"})
	ErrNotFound           = xerrors.New("not found").WithAdditionalInfo("not found", map[string]interface{}{"code": "404"})
	ErrInvalidParams      = xerrors.New("invalid params").WithExtensions(map[string]interface{}{"code": "400"})
	ErrInternal           = xerrors.New("internal error").WithExtensions(map[string]interface{}{"code": "500"})
	ErrFailedPrecondition = xerrors.New("failed precondition").WithExtensions(map[string]interface{}{"code": "412"})
	ErrAlreadyExists      = xerrors.New("already exists").WithExtensions(map[string]interface{}{"code": "409"})
	ErrResourceExhausted  = xerrors.New("resource exhausted").WithExtensions(map[string]interface{}{"code": "412"})
)

func isStorageNotFoundErr(err error) bool {
	return errors.Is(err, storage.ErrNotFound)
}

func fromStorageErr(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, storage.ErrNotFound):
		return ErrNotFound.Wrap(err)
	case errors.Is(err, storage.ErrInvalidInput):
		return ErrInvalidParams.Wrap(err)
	case errors.Is(err, storage.ErrIndexCollision):
		return ErrAlreadyExists.Wrap(err)
	case errors.Is(err, storage.ErrEntityInUse):
		return ErrAlreadyExists.Wrap(err)
	case errors.Is(err, storage.ErrInvalidConnection):
		return ErrInternal.Wrap(err)
	default:
		return xerrors.New("storage error").Consume(err)
	}
}
