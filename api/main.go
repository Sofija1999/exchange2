// Starts the CLI application, executes any given command, and shuts down
// Added to avoid having to install the egw code on the machine to execute commands
package main

import (
	"os"

	"github.com/Bloxico/exchange-gateway/api/cmd"
	"github.com/Bloxico/exchange-gateway/sofija/app"
	"github.com/urfave/cli"
)

func RunCLIApplication() {

	egwApp := app.MustInitializeApp()

	cliApp := cli.NewApp()
	cliApp.Name = "exchange-gateway"
	cliApp.Description = "Command line utility for egw development"
	cliApp.Commands = []cli.Command{
		cmd.NewDbCmd(egwApp),
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		egwApp.Logger.Error("failed running command", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func main() {
	RunCLIApplication()
}
