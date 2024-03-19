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
	configPath, configPathSet := os.LookupEnv("KEEXP_CONFIG_PATH")
	if configPathSet {
		return path.Clean(configPath), nil
	}

	xdgConfigHome, xdgConfigHomeSet := os.LookupEnv("XDG_CONFIG_HOME")
	if xdgConfigHomeSet {
		return path.Join(xdgConfigHome, "keexp/config.json"), nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(userHome, ".config/keexp/config.json"), nil
}
