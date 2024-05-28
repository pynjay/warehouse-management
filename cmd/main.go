package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"warehouse/cmd/commands"
	_ "warehouse/cmd/commands/migrate"
	"warehouse/cmd/factory"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "warehouse",
		Action: func(c *cli.Context) error {
			_, cleaner, err := factory.InitializeService(context.TODO())
			if err != nil {
				return fmt.Errorf("error initialize service. %w", err)
			}
			defer cleaner()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh

			return nil
		},
		Commands: commands.Commands,
	}

	if err := app.Run(os.Args); err != nil {
		panic(fmt.Errorf("error start app. %w", err))
	}
}
