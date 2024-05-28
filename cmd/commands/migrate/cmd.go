package migrate

import (
	"context"
	"fmt"
	"warehouse/cmd/commands"
	"warehouse/cmd/commands/migrate/migrations"
	_ "warehouse/cmd/commands/migrate/migrations/migrations"
	"warehouse/cmd/factory"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

func init() {
	commands.RegisterCommandFactory(&cli.Command{
		Name: "migration",
		Action: func(c *cli.Context) error {
			_, shutdown, err := factory.InitializeMigrationContainer()
			defer func() {
				shutdown()
			}()
			if err != nil {
				return fmt.Errorf("error initialize default container. %w", err)
			}
			goose.SetBaseFS(migrations.Migrations)
			goose.SetDialect("postgres")

			return goose.RunContext(
				context.TODO(),
				c.Args().First(),
				factory.DefaultMigrationContainer.DB(),
				"migrations",
				c.Args().Slice()[1:]...,
			)
		},
	})
}
