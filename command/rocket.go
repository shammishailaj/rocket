package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RocketCmd = &cobra.Command{
	Use:   "rocket",
	Short: "Deploy software as fast and easily as possible",
	Long:  "Deploy software as fast and easily as possible. See https://github.com/z0mbie42/rocket",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world")
	},
}
