package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RocketCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init rocket by creating a .rocket.(toml|json) configuration file",
	Long:  "Init rocket by creating a .rocket.(toml|json) configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Init....")
	},
}
