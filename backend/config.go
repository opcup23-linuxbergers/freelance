package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DB struct {
		Backend  string `json:"type"`
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
	} `json:"db"`

	Server struct {
		DevMode bool   `json:"dev_mode"`
		Port    string `json:"port"`
		Secret  string `json:"secret"`
	} `json:"server"`

	DSN string
}

func (conf *Config) composeDSN() error {
	switch conf.DB.Backend {
	case "postgres":
		conf.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432", conf.DB.Host, conf.DB.User, conf.DB.Password, conf.DB.Name)
	case "mysql":
		conf.DSN = fmt.Sprintf("%s@%s@tcp(%s)/%s", conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Name)
	case "sqlite":
		conf.DSN = conf.DB.Name
	default:
		return fmt.Errorf("unknown backend type: %s", conf.DB.Backend)
	}

	return nil
}

func configError(err error) error {
	return fmt.Errorf("config: %s", err)
}

func (conf *Config) Init() error {
	payload, err := os.ReadFile("./config.json")
	if err != nil {
		return configError(err)
	}

	if err := json.Unmarshal(payload, &conf); err != nil {
		return configError(err)
	}

	if err := conf.composeDSN(); err != nil {
		return configError(err)
	}

	return nil
}
