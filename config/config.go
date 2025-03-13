package config

import (
	"os"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	DB DB
}

func LoadConfig() *Config {
	cfg := &Config{
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}

	return cfg
}
