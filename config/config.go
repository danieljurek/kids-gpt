package config

import (
	"embed"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed default.yml
var defaultYaml embed.FS

var configIsInitialized = false
var config Config

type ErrorConfig struct {
	GeneralErrorMessage string `yaml:"generalErrorMessage"`
	ModerationMessage   string `yaml:"moderationMessage"`
	ErrorVoice          string `yaml:"errorVoice"`
}

type Config struct {
	InitialPrompt  string      `yaml:"initialPrompt"`
	SpinnerCharset int         `yaml:"spinnerCharset"`
	StopSequence   string      `yaml:"stopSequence"`
	UserName       string      `yaml:"userName"`
	GptName        string      `yaml:"gptName"`
	Model          string      `yaml:"model"`
	Errors         ErrorConfig `yaml:"errors"`
	Speed          int         `yaml:"speed"`
}

func initConfig() error {
	if configIsInitialized {
		return nil
	}

	var yamlFileContent []byte
	var err error
	if len(os.Args) >= 2 {
		if yamlFileContent, err = os.ReadFile(os.Args[1]); err != nil {
			return err
		}
	} else {
		if yamlFileContent, err = defaultYaml.ReadFile("default.yml"); err != nil {
			return err
		}
	}

	return yaml.Unmarshal(yamlFileContent, &config)
}

func GetConfig() (*Config, error) {
	if err := initConfig(); err != nil {
		return nil, err
	}

	return &config, nil
}
