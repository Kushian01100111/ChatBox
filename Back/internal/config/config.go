package config

import "os"

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port: getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, def string) string {
	res := os.Getenv(key)
	if res == "" {
		return def
	}
	return res
}
