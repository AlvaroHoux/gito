package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type GitoConfig struct {
	Model string `json:"model"`
}

func LoadConfig() (*GitoConfig, error) {
	gitoPath, err := getPathConfig()
	if err != nil {
		return nil, err
	}
	gitoConfigData, err := os.ReadFile(gitoPath)

	if err := os.IsNotExist(err); err == true {
		defaultConfig, err := createConfig(gitoPath)
		if err != nil {
			return nil, err
		}
		return defaultConfig, nil
	}

	var gitoCofig GitoConfig
	if err := json.Unmarshal(gitoConfigData, &gitoCofig); err != nil {
		return nil, err
	}
	return &gitoCofig, nil
}

func SaveConfig(model string) error {
	gitoPath, err := getPathConfig()
	if err != nil {
		return err
	}

	configData, err := json.Marshal(GitoConfig{Model: model})
	if err != nil {
		return nil
	}

	if err := os.WriteFile(gitoPath, []byte(configData), 0644); err != nil {
		return err
	}
	return nil
}

func getPathConfig() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userConfigDir, "gito", "config.json"), nil
}

func createConfig(path string) (*GitoConfig, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	defaultConfig := GitoConfig{
		Model: "granite3.3:2b",
	}

	defaultConfigData, err := json.Marshal(defaultConfig)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, defaultConfigData, 0644); err != nil {
		return nil, err
	}
	return &defaultConfig, nil
}
