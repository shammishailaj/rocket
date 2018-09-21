package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/astrocorp42/rocket/config"
	"github.com/astrocorp42/astroflow-go"
	"github.com/astrocorp42/astroflow-go/log"
	"github.com/spf13/cobra"
)

var initForce bool

func init() {
	RocketCmd.AddCommand(InitCmd)
	InitCmd.Flags().BoolVar(&initForce, "force", false, fmt.Sprintf("Force and override an existing %s.(toml|json) file", config.DefaultConfigurationFileName))
}

// InitCmd is the rocket's `init` command. It creates a configuration with default configuration
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: fmt.Sprintf("Init rocket by creating a %s configuration file", config.DefaultConfigurationFileName),
	Long:  fmt.Sprintf("Init rocket by creating a %s configuration file", config.DefaultConfigurationFileName),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := config.FindConfigFile("")
		var err error

		if debug {
			log.Config(astroflow.SetLevel(astroflow.DebugLevel))
		}

		if configFile != "" && initForce == false {
			log.Fatal(fmt.Sprintf("A configuration file already exists (%s), use --force to override", configFile))
		}

		conf := config.Default()
		filePath := config.DefaultConfigurationFileName
		buf := new(bytes.Buffer)
		err = toml.NewEncoder(buf).Encode(conf)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}
