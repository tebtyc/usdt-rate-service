package config

import (
	"fmt"
	"os"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type GRPC struct {
	Port string
}

type Config struct {
	DB         DB
	GRPC       GRPC
	Prometheus Prometheus
}

type Prometheus struct {
	Port string
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
		GRPC: GRPC{
			Port: os.Getenv("GRPC_PORT"),
		},
		Prometheus: Prometheus{
			Port: os.Getenv("PROMETHEUS_PORT"),
		},
	}

	return cfg
}

func (c *Config) GetDBURL() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name, c.DB.Port)
}
