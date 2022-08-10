package config

import "os"

type ApiConfig struct {
	ApiPort string
	ApiHost string
}
type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}
type Config struct {
	ApiConfig
	DbConfig
}

func (c Config) readConfigFile() Config {
	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	c.ApiConfig = ApiConfig{
		ApiPort: os.Getenv("API_PORT"),
		ApiHost: os.Getenv("API_HOST"),
	}
	return c
}

// DB_HOST=127.0.0.1 DB_PORT=5432 DB_DB=contohDb DB_USER=postgres DB_PASSWORD=12345678 API_HOST=localhost API_PORT=8888 go run .
func NewConfig() Config {
	cfg := Config{}
	return cfg.readConfigFile()
}
