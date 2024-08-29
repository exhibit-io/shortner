package config

import "github.com/joho/godotenv"

// Config represents the main configuration struct containing all service configs.
type Config struct {
	Redis      RedisConfig
	Redirector RedirectorConfig
	Cors       CorsConfig
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() *Config {
	godotenv.Load()
	return &Config{
		Redis:      LoadRedisConfig(),
		Redirector: LoadRedirectorConfig(),
		Cors:       LoadCorsConfig(),
	}
}
