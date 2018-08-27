package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// DefaultConfigurationFileName is the default configuration file name, without extension
const DefaultConfigurationFileName = ".rocket"

type Config struct {
	Description string `json:"description" toml:"description"`

	// providers
	Script *ScriptConfig `json:"script,omitempty" toml:"script,omitempty"`
}

// ScriptConfig is the configration for the script provider
type ScriptConfig []string

func parseConfig(configFilePath string) (Config, error) {
	var ret Config
	ext := filepath.Ext(configFilePath)
	var err error

	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return ret, err
	}

	switch ext {
	case ".toml":
		_, err = toml.Decode(string(file), &ret)
	case ".json":
		err = json.Unmarshal(file, &ret)
	default:
		err = errors.New(ext + " is not a valid configuration file extension")
	}
	if err != nil {
		return ret, err
	}

	return ret, nil
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

	ret.Description = "This is a configuration file for rocket: Deploy software as fast and easily as possible. " +
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

	tomlFile := DefaultConfigurationFileName + ".toml"
	jsonFile := DefaultConfigurationFileName + ".json"

	if fileExists(tomlFile) {
		return tomlFile
	} else if fileExists(jsonFile) {
		return jsonFile
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
			return config, fmt.Errorf("%s(.toml|json) configuration file not found. Please run \"rocket init\"", DefaultConfigurationFileName)
		}
		return config, fmt.Errorf("%s file not found.", file)
	}

	config, err = parseConfig(configFilePath)
	if err != nil {
		return config, err
	}

	return config, err
}
