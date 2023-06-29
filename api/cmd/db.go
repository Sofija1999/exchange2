package cmd

import (
	egw "github.com/Bloxico/exchange-gateway/sofija/app"
	"github.com/Bloxico/exchange-gateway/sofija/database"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func NewDbCmd(app *egw.App) cli.Command {
	return cli.Command{
		Name:    "database",
		Aliases: []string{"db"},
		Usage:   "database related actions",
		Subcommands: []cli.Command{
			NewMigrateCmd(app),
			NewResetCmd(app),
		},
	}
}

func NewMigrateCmd(app *egw.App) cli.Command {
	return cli.Command{
		Name:  "migrate",
		Usage: "execute migrations located in the ./migrations folder",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "reset",
				Usage: "before executing migrations, reset the database",
			},
			cli.StringFlag{
				Name:  "dir",
				Usage: "directory containing migrations to execute",
				Value: "migrations",
			},
		},
		Action: func(c *cli.Context) error {
			migration := database.NewMigrationProcess(app.DB, app.Logger)

			reset := c.Bool("reset")
			if reset {
				if err := performReset(migration, app); err != nil {
					return err
				}
			}

			dir := c.String("dir")

			err := migration.Migrate(dir)
			if err != nil {
				return errors.Wrap(err, "migrate db")
			}

			return nil
		},
	}
}

func NewResetCmd(app *egw.App) cli.Command {
	return cli.Command{
		Name:  "reset",
		Usage: "truncate connected database",
		Action: func(c *cli.Context) error {
			migration := database.NewMigrationProcess(app.DB, app.Logger)
			return performReset(migration, app)
		},
	}
}

func performReset(migration *database.MigrationProcess, app *egw.App) error {
	err := migration.DropSchema(app.Config.Database.Schema)
	if err != nil {
		return errors.Wrap(err, "reset db")
	}
	return nil
}
