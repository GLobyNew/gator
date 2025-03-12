package config

import (
	"encoding/json"
	"os"
)

const (
	configFileName = "/.gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := userHomeDir + configFileName
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", err
	}

	return configPath, nil

}

func Read() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var userConfig Config
	json.NewDecoder(file).Decode(&userConfig)

	return userConfig, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(c); err != nil {
		return err
	}

	return nil
}
