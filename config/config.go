package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/bloom42/astroflow-go/log"
	"github.com/bloom42/san-go"
)

// DefaultConfigurationFileName is the default configuration file name, without extension
const DefaultConfigurationFileName = ".rocket.san"

var PredefinedEnv = []string{
	"ROCKET_COMMIT_HASH",
	"ROCKET_LAST_TAG",
	"ROCKET_GIT_REPO",
}

type Config struct {
	Description string            `json:"description" san:"description"`
	Env         map[string]string `json:"env" san:"env"`

	// providers
	Script         ScriptConfig          `json:"script,omitempty" san:"script,omitempty"`
	Heroku         *HerokuConfig         `json:"heroku,omitempty" san:"heroku,omitempty"`
	GitHubReleases *GitHubReleasesConfig `json:"github_releases,omitempty" san:"github_releases,omitempty"`
	Docker         *DockerConfig         `json:"docker" san:"docker"`
	AWSS3          *AWSS3Config          `json:"aws_s3" san:"aws_s3"`
	ZeitNow        *ZeitNowConfig        `json:"zeit_now" san:"zeit_now"`
	AWSEB          *AWSEBConfig          `json:"aws_eb" san:"aws_eb"`
}

// ScriptConfig is the configration for the script provider
type ScriptConfig []string

// HerokuConfig is the configuration for the `heroku` provider
type HerokuConfig struct {
	APIKey    *string `json:"api_key" san:"api_key"`
	App       *string `json:"app" san:"app"`
	Directory *string `json:"directory" san:"directory"`
	Version   *string `json:"version" san:"version"`
}

// GitHubReleasesConfig is the configuration for the `github_releases` provider
type GitHubReleasesConfig struct {
	Name       *string  `json:"name" san:"name"`
	Body       *string  `json:"body" san:"body"`
	Prerelease *bool    `json:"prerelease" san:"prerelease"`
	Repo       *string  `json:"repo" san:"repo"`
	APIKey     *string  `json:"api_key" san:"api_key"`
	Assets     []string `json:"assets" san:"assets"`
	Tag        *string  `json:"tag" san:"tag"`
	BaseURL    *string  `json:"base_url" san:"base_url"`
	UploadURL  *string  `json:"upload_url" san:"upload_url"`
}

// DockerConfig is the configration for the docker provider
type DockerConfig struct {
	Username *string  `json:"username" san:"username"`
	Password *string  `josn:"password" san:"password"`
	Login    *bool    `json:"login" san:"login"`
	Images   []string `json:"images" san:"images"`
}

// AWSS3Config is the configration for the aws_s3 provider
type AWSS3Config struct {
	AccessKeyID     *string `json:"access_key_id" san:"access_key_id"`
	SecretAccessKey *string `json:"secret_access_key" san:"secret_access_key"`
	Region          *string `json:"region" san:"region"`
	Bucket          *string `json:"bucket" san:"bucket"`
	LocalDirectory  *string `json:"local_directory" san:"local_directory"`
	RemoteDirectory *string `json:"remote_directory" san:"remote_directory"`
}

// ZeitNowConfig is the configration for the `zeit_now` provider
type ZeitNowConfig struct {
	Token           *string           `json:"token" san:"token"`
	Directory       *string           `json:"directory" san:"directory"`
	Env             map[string]string `json:"env" san:"env"`
	Public          *bool             `json:"public" san:"public"`
	DeploymentType  *string           `json:"deployment_type" san:"deployment_type"`
	Name            *string           `json:"name" san:"name"`
	ForceNew        *bool             `json:"force_new" san:"force_new"`
	Engines         map[string]string `json:"engines" san:"engines"`
	SessionAffinity *string           `json:"session_affinity" san:"session_affinity"`
}

// AWSEBConfig is the configration for the `aws_eb` provider
type AWSEBConfig struct {
	AccessKeyID     *string `json:"access_key_id" san:"access_key_id"`
	SecretAccessKey *string `json:"secret_access_key" san:"secret_access_key"`
	Region          *string `json:"region" san:"region"`
	Application     *string `json:"application" san:"application"`
	Environment     *string `json:"environment" san:"environment"`
	S3Bucket        *string `json:"s3_bucket" san:"s3_bucket"`
	Version         *string `json:"version" san:"version"`
	Directory       *string `json:"directory" san:"directory"`
	S3Key           *string `json:"s3_key" san:"s3_key"`
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

	err = san.Unmarshal(file, &ret)

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
		"See https://github.com/bloom42/rocket"

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
