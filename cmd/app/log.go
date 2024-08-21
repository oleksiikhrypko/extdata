package main

import (
	"context"

	"github.com/Slyngshot-Team/packages/log"
)

const logConfigTag = "log"

func initDefaultLogger(ctx context.Context, serviceName string) {
	logger, err := log.BuildLoggerFromConfig(logConfigTag, serviceName, Version)
	if err != nil {
		log.Error(ctx, err, "failed to build logger from config")
		return
	}

	log.SetDefaultLogger(logger)
}
