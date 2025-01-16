package config

import (
	"log"
	"os"
)

type Config struct {
	ServerHost string
	ServerPort string
	MongoURI   string
}

func GetConfig() *Config {
	return &Config{
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "3003"),
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
	}
}

func getEnv(key, df string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Using default value for %s (%s)", key, df)
		return df
	}
	return val
}
