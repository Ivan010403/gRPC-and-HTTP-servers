package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	GRPC_server `yaml:"grpc_server"`
	DataBase    `yaml:"database"`
}

type GRPC_server struct {
	Port             int `yaml:"port" env-default:"4545"`
	MaxReadWriteConn int `yaml:"maxReadWriteConn" env-default:"10"`
	MaxCheckConn     int `yaml:"maxCheckConn" env-default:"100"`
}

type DataBase struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"8089"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"vany2003"`
	Dbname   string `yaml:"dbname" env-default:"Notepad"`
}

func ReadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is wrong")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Config file doesn't exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("can not read config %w", err)
	}

	return &cfg, nil
}
