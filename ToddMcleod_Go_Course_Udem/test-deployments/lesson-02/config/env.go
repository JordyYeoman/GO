package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DbUsername      string
	DbPassword      string
	DbConnectionUrl string
	DbPort          string
	DbName          string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load config")
		return Config{}
	}

	return Config{
		DbUsername:      getEnv("DbUsername", ""),
		DbPassword:      getEnv("DbPassword", ""),
		DbConnectionUrl: getEnv("DbConnectionUrl", ""),
		DbPort:          getEnv("DbPort", ""),
		DbName:          getEnv("DbName", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
