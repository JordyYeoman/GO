package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	OddsApiKey     string
	OddsApiBookies string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load config")
		return Config{}
	}

	// Setup from .env file
	return Config{
		OddsApiKey:     getEnv("ODDS_API_KEY", ""),
		OddsApiBookies: getEnv("ODDS_API_BOOKIES", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
