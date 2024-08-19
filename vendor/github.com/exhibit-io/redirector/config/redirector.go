package config

import "fmt"

type RedirectorConfig struct {
	Schema string
	Host   string
	Port   string
}

func (t RedirectorConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)
}

func (t RedirectorConfig) GetURI() string {
	return fmt.Sprintf("%s://%s:%s", t.Schema, t.Host, t.Port)
}

func LoadRedirectorConfig() RedirectorConfig {
	return RedirectorConfig{
		Schema: getEnv("SCHEMA", "http"),
		Host:   getEnv("HOST", "localhost"),
		Port:   getEnv("PORT", "8080"),
	}
}
