package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Archive_path string             `json:"archive_path"`
	DB_file_path string             `json:"db_file_path"`
	Logfile_path string             `json:"logfile_path"`
	Sources      []SourceConfig     `json:"sources"`
	Repositories []RepositoryConfig `json:"repositories"`
}

type SourceConfig struct {
	Active      bool              `json:"active"`
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	Screen_name string            `json:"screen_name"`
	Country     string            `json:"country"`
	Language    string            `json:"language"`
	Mediatype   string            `json:"mediatype"`
	Feed        map[string]string `json:"feed"`
}

type RepositoryConfig struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Active  bool   `json:"active"`
}

func Read(filename string) (*Config, error) {
	cfg := new(Config)
	cfg_file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cfg_file, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
