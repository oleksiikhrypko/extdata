package main

import (
	"context"

	"github.com/Slyngshot-Team/packages/log"
	psql "github.com/Slyngshot-Team/packages/storage/psql"
	"github.com/alecthomas/kingpin/v2"
)

func addMigrateCmd(ctx context.Context, app *kingpin.Application) {
	initDefaultLogger(ctx, serviceName+"::migrate")
	migrateCmd := app.Command("migrate", "Starts migration")
	migrateCommand := migrateCmd.Arg("command", "Migration command [up, up-one, down, down-one]").Required().String()
	migrateLocation := migrateCmd.Flag("location", "Migrations files location").Default("file://migrations").String()

	migrateCmd.Action(func(_ *kingpin.ParseContext) error {
		return startMigration(ctx, *migrateCommand, *migrateLocation)
	})
}

func startMigration(ctx context.Context, migrateCommand, migrateLocation string) (err error) {
	log.Info(ctx, "run migrate", "command", migrateCommand, "location", migrateLocation)
	var cmd psql.MigrateCommand
	if cmd, err = cmd.Parse(migrateCommand); err != nil {
		log.Error(ctx, err, "failed on parse migrate command")
		return err
	}

	if err = psql.StartMigrateCmd("psql", migrateLocation, cmd); err != nil {
		log.Error(ctx, err, "failed on run migrate command")
		return err
	}

	return nil
}
