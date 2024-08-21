package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/Slyngshot-Team/packages/config"
	"github.com/Slyngshot-Team/packages/log"
	srvutil "github.com/Slyngshot-Team/packages/service"
	"github.com/alecthomas/kingpin/v2"
)

var Version = "Not set" //nolint:gochecknoglobals

const serviceName = "ext_data_domain"

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "idea domain service").UsageWriter(os.Stdout)
	app.Version(Version)
	app.HelpFlag.Short('h')

	ctx, cancel := context.WithCancel(context.Background())
	srvutil.SetupGracefulShutdown(ctx, cancel)
	defer cancel()

	log.Init(log.LogFormatText, slog.LevelInfo, serviceName)

	configOnce := sync.Once{}
	var configPath *string
	var cfgErr error
	configAction := func(_ *kingpin.ParseContext) error {
		configOnce.Do(func() {
			log.Info(ctx, "load config...", "configPath", *configPath)
			cfgErr = config.LoadConfig(*configPath)
		})
		return cfgErr
	}
	configPath = app.Flag("config-file", "Service configuration file").
		Short('c').
		Action(configAction).
		Required().
		ExistingFile()

	addMigrateCmd(ctx, app)

	startCmd := app.Command("start", "Starts service")
	addWebapiCmd(ctx, startCmd)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
