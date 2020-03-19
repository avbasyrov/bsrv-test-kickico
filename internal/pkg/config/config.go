package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ABI        abi.ABI `yaml:"-"`
	ABIContent []byte  `yaml:"-"`
	ABIFile    string  `yaml:"abiFile"`
	Address    string  `yaml:"address"`
	Infura     struct {
		ProjectID  string `yaml:"project_id"`
		PrivateKey string `yaml:"private_key"`
	} `yaml:"infura"`
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
	ABIContent, err := ioutil.ReadFile(abiFilename)
	if err != nil {
		return config, err
	}
	config.ABIContent = ABIContent

	ABI, err := abi.JSON(strings.NewReader(string(ABIContent)))
	if err != nil {
		return config, err
	}
	config.ABI = ABI

	return config, err
}
