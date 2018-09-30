package main

import (
	"github.com/bloom42/rocket/commands"
	rlog "github.com/bloom42/rocket/log"
	"github.com/bloom42/astroflow-go"
	"github.com/bloom42/astroflow-go/log"
)

func main() {
	log.Config(
		astroflow.SetFormatter(rlog.NewCLIFormatter()),
		astroflow.SetLevel(astroflow.InfoLevel),
	)
	if err := commands.RocketCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
