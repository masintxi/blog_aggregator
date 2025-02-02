package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

/*
>> to migrate, from the sql/schema dir:
goose postgres <connection_string> up (or down)
>> where <connection_string> is:
postgres://postgres:postgres@localhost:5432/gator
>> the whole thing:
goose postgres postgres://postgres:postgres@localhost:5432/gator up
*/

func getConfigPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homePath, configFileName), nil
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

	var data Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return Config{}, err
	}
	return data, nil
}

func write(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	return write(cfg)
}
