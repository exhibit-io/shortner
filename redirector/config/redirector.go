package config

import "fmt"

type RedirectorConfig struct {
	Schema    string
	Host      string
	Port      string
	PublicURL string
}

func (t RedirectorConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)
}

func (t RedirectorConfig) GetURI() string {
	return fmt.Sprintf("%s://%s:%s", t.Schema, t.Host, t.Port)
}

func LoadRedirectorConfig() RedirectorConfig {
	config := RedirectorConfig{
		Schema: getEnv("SCHEMA", "http"),
		Host:   getEnv("HOST", "localhost"),
		Port:   getEnv("PORT", "8080"),
	}

	// Set the PublicURL field after the object is initialized
	config.PublicURL = getEnv("PUBLIC_URL", config.GetURI())

	return config
}
