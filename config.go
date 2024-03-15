package main

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

func readConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	configJson, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configJson, &config)
	return &config, err
}

func getConfigPath() (string, error) {
	configBasePath, configHomeSet := os.LookupEnv("XDG_CONFIG_HOME")
	if !configHomeSet {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configBasePath = path.Join(userHome, ".config")
	}

	return path.Join(configBasePath, "keexp/config.json"), nil
}
