package commands

import (
	"fmt"
	"os"

	"github.com/astrocorp42/rocket/config"
	"github.com/astrocorp42/rocket/providers/script"
	"github.com/spf13/cobra"
)

// RocketCmd is the rocket's root command. It's used to actually deploy
var RocketCmd = &cobra.Command{
	Use:   "rocket",
	Short: "Deploy software as fast and easily as possible",
	Long:  "Deploy software as fast and easily as possible. See https://github.com/z0mbie42/rocket",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Get()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		script.Deploy(conf)
	},
}
