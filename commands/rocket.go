package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astrocorp42/rocket/config"
	"github.com/astrocorp42/rocket/providers/docker"
	"github.com/astrocorp42/rocket/providers/ghreleases"
	"github.com/astrocorp42/rocket/providers/heroku"
	"github.com/astrocorp42/rocket/providers/script"
	"github.com/astroflow/astroflow-go"
	"github.com/astroflow/astroflow-go/log"
	"github.com/spf13/cobra"
)

var rocketConfigPath string
var debug bool

func init() {
	RocketCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Display debug information")
	RocketCmd.Flags().StringVarP(&rocketConfigPath, "config", "c", "", "Use the specified configuration file (and set it's directory as the working directory")
}

// RocketCmd is the rocket's root command. It's used to actually deploy
var RocketCmd = &cobra.Command{
	Use:   "rocket",
	Short: "Automated software delivery as fast and easy as possible",
	Long:  "Automated software delivery as fast and easy as possible. rocket is the D in CI/CD. See https://github.com/astrocorp42/rocket",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if debug {
			log.Config(astroflow.SetLevel(astroflow.DebugLevel))
		}

		// change working directory as the file's
		if rocketConfigPath != "" {
			dir := filepath.Dir(rocketConfigPath)
			err = os.Chdir(dir)
			if err != nil {
				log.Fatal(err.Error())
			}
			rocketConfigPath = filepath.Base(rocketConfigPath)
		}

		conf, err := config.Get(rocketConfigPath)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.With("configuration", conf).Debug("")
		log.With("env", os.Environ()).Debug("")

		// script
		if conf.Script != nil {
			log.Debug("script: starting provider")
			err = script.Deploy(conf.Script)
			if err != nil {
				log.Fatal(fmt.Sprintf("script: %v", err))
			}
		} else {
			log.Debug("script: provider is empty")
		}

		// heroku
		if conf.Heroku != nil {
			log.Debug("heroku: starting provider")
			err = heroku.Deploy(*conf.Heroku)
			if err != nil {
				log.Fatal(fmt.Sprintf("heroku: %v", err))
			}
		} else {
			log.Debug("heroku: provider is empty")
		}

		// github_releases
		if conf.GitHubReleases != nil {
			log.Debug("github_releases: starting provider")
			err = ghreleases.Deploy(*conf.GitHubReleases)
			if err != nil {
				log.Fatal(fmt.Sprintf("github_releases: %v", err))
			}
		} else {
			log.Debug("github_releases: provider is empty")
		}

		// docker
		if conf.Docker != nil {
			log.Debug("docker: starting provider")
			err = docker.Deploy(*conf.Docker)
			if err != nil {
				log.Fatal(fmt.Sprintf("docker: %v", err))
			}
		} else {
			log.Debug("docker: provider is empty")
		}
	},
}
