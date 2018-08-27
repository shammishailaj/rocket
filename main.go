package main

import (
	"fmt"
	"os"

	"github.com/astrocorp42/rocket/command"
)

func main() {
	if err := command.RocketCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
