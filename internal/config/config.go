package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(home, CONFIG_FILE_NAME)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
