package webapi

import (
	"fmt"
	"net/http"

	"ext-data-domain/internal/server/webapi/resolver/graph"
	"ext-data-domain/internal/server/webapi/resolver/openapi"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/slyngshot-al/packages/log"
	"github.com/slyngshot-al/packages/service"
	// echotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

const ConfigTag = service.WebServerConfigTag

type Middleware func(handler http.Handler) http.Handler

func New(conf service.WebServerConfig, apiConf *openapi.Config, graphConf *graph.Config, checkFns ...service.CheckFn) (*service.WebServer, error) {
	srv, err := service.NewWebServer(conf, service.GetEchoHealthHandler(checkFns...), log.GetLogger("webserver"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize web server: %w", err)
	}

	// recover
	srv.Use(NewRecovery())

	// metrics
	srv.Use(echoprometheus.NewMiddleware("webapi")) // adds middleware to gather metrics

	// init meta /meta
	meta := srv.Group("/meta")
	meta.File("/policy", "./configs/policy.csv")

	// trace
	// srv.Use(echotrace.Middleware(
	// 	echotrace.WithServiceName("webapi"),
	// 	echotrace.WithIgnoreRequest(func(c echo.Context) bool {
	// 		return c.Path() == "/health"
	// 	}),
	// ))

	// init graphql
	if err := graph.Register(srv, graphConf); err != nil {
		return nil, err
	}

	// init openapi
	if err := openapi.Register(srv, apiConf); err != nil {
		return nil, err
	}

	return srv, nil
}
