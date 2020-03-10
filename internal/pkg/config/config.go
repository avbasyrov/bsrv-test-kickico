package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ABI     []byte `yaml:"-"`
	ABIFile string `yaml:"abiFile"`
}

func Read(configPath string, configName string) (Config, error) {
	var config Config

	filename, err := filepath.Abs(configPath + "/" + configName)
	if err != nil {
		return config, err
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	abiFilename, err := filepath.Abs(configPath + "/" + config.ABIFile)
	if err != nil {
		return config, err
	}
	ABI, err := ioutil.ReadFile(abiFilename)
	if err != nil {
		return config, err
	}

	config.ABI = ABI

	return config, err
}
