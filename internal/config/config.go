package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME = ".gatorconfig.json"
const CONFIG_TMP_FILE_NAME = ".gatorconfig.json.tmp"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	home, err := os.UserHomeDir()
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

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tmpFilePath := filepath.Join(home, CONFIG_TMP_FILE_NAME)

	tmp, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}

	err = json.NewEncoder(tmp).Encode(cfg)
	if err != nil {
		tmp.Close()
		return err
	}

	err = tmp.Close()
	if err != nil {
		return err
	}

	filePath := filepath.Join(home, CONFIG_FILE_NAME)
	return os.Rename(tmpFilePath, filePath)
}
