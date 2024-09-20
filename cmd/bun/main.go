package main

import (
	"os"
	"strings"

	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"github.com/yesyash/yaskin-backend/cmd/bun/migrations"
	"github.com/yesyash/yaskin-backend/internal/database"
	"github.com/yesyash/yaskin-backend/internal/logger"
)

func main() {
	app := &cli.App{
		Name: "bun",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			newDBCommand(migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)
					group, err := migrator.Migrate(c.Context)

					if err != nil {
						return err
					}

					if group.ID == 0 {
						logger.Info("there are no new migrations to run\n")
						return nil
					}

					logger.Info("migrated to ", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)
					group, err := migrator.Rollback(c.Context)

					if err != nil {
						return err
					}

					if group.ID == 0 {
						logger.Info("there are no groups to roll back\n")
						return nil
					}

					logger.Info("rolled back ", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					return migrator.Lock(c.Context)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)

					return migrator.Unlock(c.Context)
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)
					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(c.Context, name)

					if err != nil {
						return err
					}

					for _, mf := range files {
						logger.Info("created migration ", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)
					ms, err := migrator.MigrationsWithStatus(c.Context)

					if err != nil {
						return err
					}

					logger.Info("migrations: ", ms)
					logger.Info("unapplied migrations: ", ms.Unapplied())
					logger.Info("last migration group: ", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					db := database.New()
					defer db.Close()

					migrator := migrate.NewMigrator(db, migrations)
					group, err := migrator.Migrate(c.Context, migrate.WithNopMigration())

					if err != nil {
						return err
					}

					if group.ID == 0 {
						logger.Info("there are no new migrations to mark as applied\n")
						return nil
					}

					logger.Info("marked as applied ", group)
					return nil
				},
			},
		},
	}
}
