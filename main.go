package main

import (
	"github.com/bloom42/astro-go"
	"github.com/bloom42/astro-go/log"
	"github.com/bloom42/rocket/commands"
	rlog "github.com/bloom42/rocket/log"
)

func main() {
	log.Config(
		astro.SetFormatter(rlog.NewCLIFormatter()),
		astro.SetLevel(astro.InfoLevel),
	)
	if err := commands.RocketCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
