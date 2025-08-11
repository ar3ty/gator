package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed attempt to get an homedir path: %w", err)
	}

	pathToConfig := filepath.Join(homePath, configFileName)
	return pathToConfig, nil
}

func Read() (Config, error) {
	config := Config{}

	path, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("cannot get config path: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed reading file: %w", err)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, fmt.Errorf("failed unmarshaling json: %w", err)
	}

	return config, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("cannot get config path: %w", err)
	}

	configJson, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed marshaling json: %w", err)
	}

	err = os.WriteFile(path, configJson, 0644)
	if err != nil {
		return fmt.Errorf("error writing in configfile: %w", err)
	}

	return nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}
