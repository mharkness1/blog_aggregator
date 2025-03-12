package config

import (
	"encoding/json"
	"io"
	"os"
)

type ConfigStruct struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const configPath string = ".gatorconfig.json"

func ReadConfig() (ConfigStruct, error) {
	var config ConfigStruct

	file, err := os.Open(configPath)
	if err != nil {
		return ConfigStruct{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return ConfigStruct{}, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return ConfigStruct{}, err
	}
	return config, nil
}

func (ConfigStruct) SetUser() error {
	return nil
}
