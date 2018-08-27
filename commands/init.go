package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/astrocorp42/rocket/config"
	"github.com/spf13/cobra"
)

var initFormat string
var initForce bool

func init() {
	RocketCmd.AddCommand(InitCmd)
	InitCmd.Flags().StringVarP(&initFormat, "format", "f", "toml", "Format of the configuration file. Valid values are [toml, json]")
	InitCmd.Flags().BoolVar(&initForce, "force", false, fmt.Sprintf("Force and override an existing %s.(toml|json) file", config.DefaultConfigurationFileName))
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init rocket by creating a .rocket.(toml|json) configuration file",
	Long:  "Init rocket by creating a .rocket.(toml|json) configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := config.FindConfigFile()
		var err error

		if configFile != "" && initForce == false {
			fmt.Fprintf(os.Stderr, "A configuration file already exists (%s), use --force to override\n", configFile)
			os.Exit(1)
		}

		conf := config.Default()
		filePath := config.DefaultConfigurationFileName
		buf := new(bytes.Buffer)

		switch initFormat {
		case "toml":
			err = toml.NewEncoder(buf).Encode(conf)
			filePath += ".toml"
		case "json":
			err = json.NewEncoder(buf).Encode(conf)
			filePath += ".json"
		default:
			err = fmt.Errorf("%s is not a valid configuration file format", initFormat)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
