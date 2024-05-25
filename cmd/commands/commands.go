package commands

import (
	"github.com/urfave/cli/v2"
)

var Commands []*cli.Command

func RegisterCommandFactory(c *cli.Command) {
	Commands = append(Commands, c)
}
