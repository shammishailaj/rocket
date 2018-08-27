package main

import (
	"fmt"
	"os"

	"github.com/astrocorp42/rocket/commands"
)

func main() {
	if err := commands.RocketCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
