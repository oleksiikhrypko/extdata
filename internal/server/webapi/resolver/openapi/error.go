package openapi

import (
	"errors"
	"net/http"
	"strings"

	api "ext-data-domain/internal/server/webapi/api/openapi"
	"ext-data-domain/internal/service"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/slyngshot-al/packages/xerrors"
)

func bindError(c echo.Context, err error) error {
	msg := api.Error{
		Message: err.Error(),
	}
	return c.JSON(http.StatusBadRequest, msg)
}

func handleError(_ echo.Context, err error) error {
	msg := api.Error{
		Message: err.Error(),
	}

	var xerr *xerrors.Error
	if errors.As(err, &xerr) {
		msg.Fields = xerr.Extensions()
	}

	switch {
	case errors.Is(err, service.ErrForbidden):
		return echo.NewHTTPError(http.StatusForbidden, msg)
	case errors.Is(err, service.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, msg)
	case errors.Is(err, service.ErrInvalidParams):
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	case errors.Is(err, service.ErrAlreadyExists):
		return echo.NewHTTPError(http.StatusPreconditionFailed, msg)
	case errors.Is(err, service.ErrResourceExhausted):
		return echo.NewHTTPError(http.StatusPreconditionFailed, msg)
	case errors.Is(err, service.ErrFailedPrecondition):
		return echo.NewHTTPError(http.StatusPreconditionFailed, msg)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, msg)
	}
}

func multiErrorHandler(me openapi3.MultiError) *echo.HTTPError {
	var (
		multiErr   openapi3.MultiError
		requestErr *openapi3filter.RequestError
	)

	message := api.Error{
		Fields:  make(map[string]any),
		Message: me.Error(),
	}

	for _, err := range me {
		if err == nil {
			continue
		}
		// check if error is multi error try to fill fields
		if ok := errors.As(err, &multiErr); ok {
			fillFieldsInfoFromMultiError(multiErr, message.Fields)
		}
		// check if error is request error try to fill fields
		if ok := errors.As(err, &requestErr); ok {
			fillFieldsInfoFromRequestError(requestErr, message.Fields)
		}
	}

	return &echo.HTTPError{
		Code:     http.StatusBadRequest,
		Message:  message,
		Internal: me,
	}
}

func fillFieldsInfoFromRequestError(requestErr *openapi3filter.RequestError, fields map[string]any) {
	if requestErr == nil || requestErr.Err == nil {
		return
	}

	var (
		parseErr *openapi3filter.ParseError
	)

	// check if error is parse error try to fill fields
	if ok := errors.As(requestErr.Err, &parseErr); ok {
		fillFieldsInfoFromParseError(parseErr, fields, requestErr.Parameter)
	}
}

func fillFieldsInfoFromParseError(parseErr *openapi3filter.ParseError, fields map[string]any, param *openapi3.Parameter) {
	if parseErr == nil || param == nil {
		return
	}
	fields[param.Name] = parseErr.Reason
}

func fillFieldsInfoFromMultiError(multiErr openapi3.MultiError, fields map[string]any) {
	if multiErr == nil {
		return
	}
	for _, err := range multiErr {
		if err == nil {
			continue
		}
		// check if error is schema error try to fill fields
		var schemaErr *openapi3.SchemaError
		if ok := errors.As(err, &schemaErr); ok {
			fillFieldsInfoFromSchemaError(schemaErr, fields)
		}
	}
}

func fillFieldsInfoFromSchemaError(schemaError *openapi3.SchemaError, fields map[string]any) {
	if schemaError == nil {
		return
	}
	fieldPath := strings.Join(schemaError.JSONPointer(), ".")
	fields[fieldPath] = schemaError.Reason
}
