package config

import (
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func NewWithDefaults() *Config {
	// create copy of defaults
	config := Defaults

	return &config
}

func (cfg *Config) LoadYaml(filePath string) {

	// read raw data
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Fatalf("failed to load config file: %v error: %+v", filePath, err)
	}

	// parse config and overwrite set values
	err = yaml.Unmarshal(contents, cfg)
	if err != nil {
		logrus.Fatalf("failed to parse config file: %v error: %+v: %v", filePath, err)
	}
}

func (cfg *Config) ToString() string {
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}

	return string(bytes)
}
