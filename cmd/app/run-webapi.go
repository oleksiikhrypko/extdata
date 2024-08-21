package main

import (
	"context"
	"fmt"

	"ext-data-domain/internal/server/webapi"
	apigraph "ext-data-domain/internal/server/webapi/resolver/graph"
	"ext-data-domain/internal/server/webapi/resolver/openapi"
	"ext-data-domain/internal/service"

	"github.com/Slyngshot-Team/packages/auth"
	"github.com/Slyngshot-Team/packages/config"
	"github.com/Slyngshot-Team/packages/filestorage/awss3"
	"github.com/Slyngshot-Team/packages/log"
	psrv "github.com/Slyngshot-Team/packages/service"
	"github.com/Slyngshot-Team/packages/storage/psql"
	"github.com/alecthomas/kingpin/v2"
)

func addWebapiCmd(ctx context.Context, root *kingpin.CmdClause) {
	root.Command("webapi", "Starts WEB API Server").Action(func(_ *kingpin.ParseContext) error {
		svcName := serviceName + "_webapi"
		initDefaultLogger(ctx, svcName)
		err := startWebApi(ctx, svcName)
		if err != nil {
			log.Error(ctx, err, "failed on run webapi server")
			return err
		}
		return nil
	})
}

func startWebApi(ctx context.Context, svcName string) error {
	psrv.InitTracing(svcName, Version)()

	err := service.InitMetrics(svcName, ":9090")
	if err != nil {
		log.Error(ctx, err, "failed on register metrics")
		return err
	}

	db, err := psql.InitDBFromConfig(ctx, "psql", svcName)
	if err != nil {
		log.Error(ctx, err, "failed on init db")
		return err
	}
	defer psql.GetStopDBFn(ctx, db)

	fileStorage, err := awss3.InitFromDefaultConf(ctx)
	if err != nil {
		log.Error(ctx, err, "failed on init file storage")
		return err
	}

	// init worldlogo service
	var worldlogoSvcCfg service.WorldLogoServiceConfig
	err = config.Unmarshal("idea-service", &worldlogoSvcCfg)
	if err != nil {
		log.Error(ctx, err, "failed on unmarshal 'idea-service' config")
		return err
	}
	worldlogoSvcCfg.DbConn = db
	worldlogoSvcCfg.FileStorage = fileStorage
	worldlogoService := service.NewWorldLogoService(worldlogoSvcCfg)

	// init auth
	var casCfg auth.CasbinConfig
	err = config.Unmarshal(auth.CasbinConfigTag, &casCfg)
	if err != nil {
		return fmt.Errorf("failed to parse casbin config: %w", err)
	}
	var authCfg auth.AuthnConfig
	if err = config.Unmarshal(auth.AuthnConfigTag, &authCfg); err != nil {
		return fmt.Errorf("failed to unmarshal auth config: %w", err)
	}

	auth.SetEnforcerProvider(auth.NewEnforcerWithHeaders)

	authMid, err := auth.HttpAuthInterceptor2(authCfg, &casCfg)
	if err != nil {
		log.Error(ctx, err, "failed on build auth middleware")
		return err
	}
	wsGraphInitFunc, err := auth.WebsocketInitFunc(authCfg, &casCfg, "/graph/query")
	if err != nil {
		return fmt.Errorf("failed to build websocket init func: %w", err)
	}

	// init graph config
	graphConf := &apigraph.Config{
		AuthMid:          authMid,
		InitFunc:         wsGraphInitFunc,
		WorldlogoService: worldlogoService,
	}

	// init openapi config
	apiConf := &openapi.Config{
		AuthMid:          authMid,
		WorldlogoService: worldlogoService,
	}

	// init webapi server
	var conf psrv.WebServerConfig
	err = config.Unmarshal(webapi.ConfigTag, &conf)
	if err != nil {
		log.Error(ctx, err, "failed on unmarshal 'webapi' config")
		return err
	}
	conf.ServiceName = svcName

	// build server
	server, err := webapi.New(conf, apiConf, graphConf, db.PingContext)
	if err != nil {
		log.Error(ctx, err, "failed on build web server")
		return err
	}

	// run
	runCtx := log.CtxWithValues(ctx, "service-name", svcName)
	return psrv.RunInParallel(runCtx,
		psrv.RunnableFunc(server.Run),
	)
}
