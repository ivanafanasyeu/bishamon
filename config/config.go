package config

import "os"

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

type MongoConfig struct {
	CONNECT_URI string
}

type AppConfig struct {
	PORT             string
	JWT_TOKEN_SECRET string
}

type Config struct {
	App   AppConfig
	Mongo MongoConfig
}

func New() *Config {
	return &Config{
		App: AppConfig{
			PORT:             getEnv("PORT", "4000"),
			JWT_TOKEN_SECRET: getEnv("JWT_TOKEN_SECRET", ""),
		},
		Mongo: MongoConfig{
			CONNECT_URI: getEnv("MONGODB_URI", ""),
		},
	}
}
