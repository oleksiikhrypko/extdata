package openapi

import (
	"context"
	"errors"
	"net/http"

	"ext-data-domain/internal/model"
	api "ext-data-domain/internal/server/webapi/api/openapi"

	"github.com/Slyngshot-Team/packages/storage/psql"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	echo "github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	apimid "github.com/oapi-codegen/echo-middleware"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Handler struct {
	worldlogoService WorldlogoService
}

type WorldlogoService interface {
	GetWorldLogoById(ctx context.Context, id string) (res model.WorldLogo, err error)
	GetWorldLogos(ctx context.Context, ops model.WorldLogosQueryOptions, sort []psql.Sort, pg psql.Pagination) (res []model.WorldLogo, err error)
	SaveWorldLogo(ctx context.Context, apiKey string, input model.WorldLogoInput) (id string, err error)
	DeleteWorldLogo(ctx context.Context, apiKey string, ids ...string) (err error)
}

type Config struct {
	WorldlogoService WorldlogoService
	AuthMid          Middleware
}

func checkConfig(conf *Config) error {
	if conf == nil {
		return errors.New("openapi config is empty")
	}
	if conf.WorldlogoService == nil {
		return errors.New("openapi world logo service is nil")
	}
	return nil
}

type Server interface {
	Group(prefix string, middleware ...echo.MiddlewareFunc) *echo.Group
}

type Middleware func(handler http.Handler) http.Handler

func Register(e Server, conf *Config) error {
	if err := checkConfig(conf); err != nil {
		return err
	}

	// init swagger
	swagger, err := api.GetSwagger()
	if err != nil {
		return err
	}
	swagger.Servers = openapi3.Servers{
		&openapi3.Server{
			Extensions:  nil,
			URL:         "/api/",
			Description: "",
			Variables:   nil,
		},
	}

	// init swagger /docs
	docs := e.Group("/docs")
	docs.Use(echomw.StaticWithConfig(echomw.StaticConfig{
		Root: "openapi",
	}))

	// init API Group for /api
	apiGp := e.Group("/api")
	apiGp.Use(echo.WrapMiddleware(conf.AuthMid))
	apiGp.Use(apimid.OapiRequestValidatorWithOptions(swagger, &apimid.Options{
		ErrorHandler: nil,
		Options: openapi3filter.Options{
			AuthenticationFunc: authFunc,
			MultiError:         true,
		},
		ParamDecoder:      nil,
		UserData:          nil,
		Skipper:           nil,
		MultiErrorHandler: multiErrorHandler,
	}))
	// add api handler
	handler := &Handler{
		worldlogoService: conf.WorldlogoService,
	}
	api.RegisterHandlers(apiGp, handler)

	return nil
}

func authFunc(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	return nil
}
