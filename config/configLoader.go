package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoaderConfigs(configPath string) error {
	var configTemp ConfigsValue
	file, err := os.Open(configPath)
	if err != nil {
		return err
	} else {
		decode := yaml.NewDecoder(file)
		decode.Decode(&configTemp)
		Config = configTemp
		return nil
	}
}
