package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(configFile string) (Config, error) {
	Debug("loading config file: ", configFile)

	var config = Config{}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	Debug("loaded config file: ", &config)

	return config, nil
}

func ParseMonth(digit int) (Month, error) {
	if digit < 1 || digit > 12 {
		return 0, errors.New("month must be between 1 and 12")
	}
	return Month(digit), nil
}
