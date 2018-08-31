package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/astroflow/astroflow-go/log"
)

// DefaultConfigurationFileName is the default configuration file name, without extension
const DefaultConfigurationFileName = ".rocket.toml"

var PredefinedEnv = []string{
	"ROCKET_COMMIT_HASH",
	"ROCKET_LAST_TAG",
	"ROCKET_GIT_REPO",
}

type Config struct {
	Description string            `json:"description" toml:"description"`
	Env         map[string]string `json:"env" toml:"env"`

	// providers
	Script         ScriptConfig          `json:"script,omitempty" toml:"script,omitempty"`
	Heroku         *HerokuConfig         `json:"heroku,omitempty" toml:"heroku,omitempty"`
	GitHubReleases *GitHubReleasesConfig `json:"github_releases,omitempty" toml:"github_releases,omitempty"`
	Docker         *DockerConfig         `json:"docker" toml:"docker"`
}

// ScriptConfig is the configration for the script provider
type ScriptConfig []string

type HerokuConfig struct {
	APIKey    *string `json:"api_key" toml:"api_key"`
	App       *string `json:"app" toml:"app"`
	Directory *string `json:"directory" toml:"directory"`
	Version   *string `json:"version" toml:"version"`
}

type GitHubReleasesConfig struct {
	Name       *string  `json:"name" toml:"name"`
	Body       *string  `json:"body" toml:"body"`
	Prerelease *bool    `json:"prerelease" toml:"prerelease"`
	Repo       *string  `json:"repo" toml:"repo"`
	APIKey     *string  `json:"api_key" toml:"api_key"`
	Assets     []string `json:"assets" toml:"assets"`
	Tag        *string  `json:"tag" toml:"tag"`
}

// DockerConfig is the configration for the docker provider
type DockerConfig struct {
	Username *string  `json:"username" toml:"username"`
	Password *string  `josn:"password" toml:"password"`
	Login    *bool    `json:"login" toml:"login"`
	Images   []string `json:"images" toml:"images"`
}

// ExpandEnv 'fix' os.ExpandEnv by allowing to use $$ to escape a dollar e.g: $$HOME -> $HOME
func ExpandEnv(s string) string {
	os.Setenv("ROCKET_DOLLAR", "$")
	return os.ExpandEnv(strings.Replace(s, "$$", "${ROCKET_DOLLAR}", -1))
}

func parseConfig(configFilePath string) (Config, error) {
	var ret Config
	var err error

	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return ret, err
	}

	_, err = toml.Decode(string(file), &ret)

	return ret, err
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// Default return a Config struct filled with default configuration
func Default() Config {
	var ret Config

	ret.Description = "This is a configuration file for rocket: automated software delivery as fast and easy as possible. " +
		"See https://github.com/astrocorp42/rocket"

	return ret
}

// FindConfigFile return the path of the first configuration file found
// it returns an emtpy string if none is found
func FindConfigFile(file string) string {
	if file != "" {
		if fileExists(file) {
			return file
		}
		return ""
	}

	if fileExists(DefaultConfigurationFileName) {
		return DefaultConfigurationFileName
	}

	return ""
}

// Get return the parsed found configuration file or an error
func Get(file string) (Config, error) {
	var err error
	var config Config

	configFilePath := FindConfigFile(file)

	if configFilePath == "" {
		if file == "" {
			return config, fmt.Errorf("%s configuration file not found. Please run \"rocket init\"", DefaultConfigurationFileName)
		}
		return config, fmt.Errorf("%s file not found.", file)
	}

	config, err = parseConfig(configFilePath)
	if err != nil {
		return config, err
	}

	err = setPredefinedEnv()
	if err != nil {
		return config, err
	}

	err = parseEnv(config)
	if err != nil {
		return config, err
	}

	return config, err
}

// set the default env variables
// it does not overwrite the already existing
func setPredefinedEnv() error {
	if os.Getenv("ROCKET_COMMIT_HASH") == "" {
		v := ""
		out, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err == nil {
			v = strings.TrimSpace(string(out))
		} else {
			log.With("err", err, "var", "ROCKET_COMMIT_HASH").Debug("error setting env var")
		}
		err = os.Setenv("ROCKET_COMMIT_HASH", v)
		if err != nil {
			return err
		}
	}

	if os.Getenv("ROCKET_LAST_TAG") == "" {
		v := ""
		out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
		if err == nil {
			v = strings.TrimSpace(string(out))
		} else {
			log.With("err", err, "var", "ROCKET_LAST_TAG").Debug("error setting env var")
		}
		err = os.Setenv("ROCKET_LAST_TAG", v)
		if err != nil {
			return err
		}
	}

	if os.Getenv("ROCKET_GIT_REPO") == "" {
		v := ""
		out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
		if err == nil {
			parts := strings.Split(strings.TrimSpace(string(out)), ":")
			parts = strings.Split(parts[len(parts)-1], "/")
			repo := parts[len(parts)-2] + "/" + parts[len(parts)-1]
			v = strings.Replace(repo, ".git", "", -1)
		} else {
			log.With("err", err, "var", "ROCKET_GIT_REPO").Debug("error setting env var")
		}
		err = os.Setenv("ROCKET_GIT_REPO", v)
		if err != nil {
			return err
		}
	}

	return nil
}

func isPredefined(key string) bool {
	for _, v := range PredefinedEnv {
		if v == key {
			return true
		}
	}

	return false
}

// parseVariables parse the 'variables' field of the configuration, expand them and set them as env
func parseEnv(conf Config) error {
	if conf.Env != nil {
		for key, value := range conf.Env {
			var err error
			key = strings.ToUpper(key)
			if os.Getenv(key) == "" || isPredefined(key) {
				err = os.Setenv(key, ExpandEnv(value))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
