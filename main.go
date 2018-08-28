package main

import (
	"github.com/astrocorp42/rocket/commands"
	"github.com/astroflow/astroflow-go"
	"github.com/astroflow/astroflow-go/log"
)

func main() {
	log.Config(
		astroflow.SetFormatter(astroflow.NewCLIFormatter()),
		astroflow.SetLevel(astroflow.InfoLevel),
	)
	if err := commands.RocketCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
