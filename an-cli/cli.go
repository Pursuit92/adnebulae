package main

import (
	"github.com/Pursuit92/cli"
)

type CLI struct {
	*cli.Command
}

func NewCLI() *CLI {
	return &CLI{cli.New("an-cli", "Adnebulae CLI: Manage Openstack instances along with their configurations", cli.HelpOnly)}
}
