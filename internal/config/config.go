package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type ConfigStruct struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const configName string = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configName)
	return fullPath, nil
}

func ReadConfig() (ConfigStruct, error) {
	var cfg ConfigStruct

	fullPath, err := getConfigFilePath()
	if err != nil {
		return ConfigStruct{}, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return ConfigStruct{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return ConfigStruct{}, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return ConfigStruct{}, err
	}
	return cfg, nil
}

func WriteConfig(cfg ConfigStruct) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	json, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(fullPath, json, 0222)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigStruct) SetUser(userName string) error {
	cfg.CurrentUser = userName
	return WriteConfig(*cfg)
}
