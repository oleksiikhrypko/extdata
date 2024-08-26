package graph

import (
	"context"
	"errors"
	"net/http"
	"time"

	"ext-data-domain/internal/model"
	graphapi "ext-data-domain/internal/server/webapi/api/graph"

	"github.com/99designs/gqlgen/graphql"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	gqlext "github.com/99designs/gqlgen/graphql/handler/extension"
	gqllru "github.com/99designs/gqlgen/graphql/handler/lru"
	gqltransport "github.com/99designs/gqlgen/graphql/handler/transport"
	gqlplayground "github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
	"github.com/slyngshot-al/packages/storage/psql"
	"github.com/slyngshot-al/packages/xerrors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	worldlogoService WorldlogoService
}

type WorldlogoService interface {
	GetWorldLogos(ctx context.Context, ops model.WorldLogosQueryOptions, sort []psql.Sort, pg psql.Pagination) (items []model.WorldLogo, err error)
	GetWorldLogosCount(ctx context.Context, ops model.WorldLogosQueryOptions) (count uint64, err error)
}

type Config struct {
	AuthMid          Middleware
	InitFunc         func(ctx context.Context, initPayload gqltransport.InitPayload) (context.Context, *gqltransport.InitPayload, error)
	WorldlogoService WorldlogoService
}

func checkConfig(conf *Config) error {
	if conf == nil {
		return errors.New("graph config is empty")
	}
	if conf.WorldlogoService == nil {
		return errors.New("graph world logo service is nil")
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

	graphGp := e.Group("/graph")
	// init graphql -> playground
	graphGp.Any("", echo.WrapHandler(gqlplayground.Handler("GraphQL playground", "/graph/query")))

	// init graphql query resolver
	es := graphapi.NewExecutableSchema(graphapi.Config{
		Resolvers: &Resolver{
			worldlogoService: conf.WorldlogoService,
		},
	})

	// init graphql server
	srv := gqlhandler.New(es)
	srv.AddTransport(&gqltransport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		InitFunc:              conf.InitFunc,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// return r.Host == "example.org"
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	srv.AddTransport(gqltransport.Options{})
	srv.AddTransport(gqltransport.SSE{})
	srv.AddTransport(gqltransport.GET{})
	srv.AddTransport(gqltransport.POST{})
	srv.AddTransport(gqltransport.MultipartForm{})

	srv.SetQueryCache(gqllru.New(1000))

	srv.Use(gqlext.Introspection{})
	srv.Use(gqlext.AutomaticPersistedQuery{
		Cache: gqllru.New(100),
	})

	// set error handler
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		var xerr *xerrors.Error
		if errors.As(e, &xerr) {
			err.Extensions = xerr.Extensions()
			err.Message = xerr.Message()
		}

		return err
	})

	// handle graphql query with data loader
	queryGp := graphGp.Group("/query")
	queryGp.Use(echo.WrapMiddleware(conf.AuthMid))
	queryGp.Any("", echo.WrapHandler(srv))

	return nil
}
