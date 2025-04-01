package config

import (
	"log"
	"os"
)

type Config struct {
	ServerHost     string
	ServerPort     string
	SpotURL        string
	UserURL        string
	ReservationURL string
}

func GetConfig() *Config {
	return &Config{
		ServerHost:     getEnv("SERVER_HOST", "localhost"),
		ServerPort:     getEnv("SERVER_PORT", "3004"),
		SpotURL:        getEnv("SPOT_URL", "http://localhost:3002"),
		UserURL:        getEnv("USER_URL", "http://localhost:3001"),
		ReservationURL: getEnv("RESERVATION_URL", "http://localhost:3003"),
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
