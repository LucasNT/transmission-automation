package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigNotExist error = errors.New("Config don't exist")
)

func LoaderConfigs(configPath string) error {
	var configTemp ConfigsValue
	file, err := os.Open(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return ErrConfigNotExist
	} else if err != nil {
		return err
	}
	decode := yaml.NewDecoder(file)
	err = decode.Decode(&configTemp)
	if err != nil {
		return fmt.Errorf("Failed to read the config file: %w", err)
	}
	Config = configTemp
	return nil
}
