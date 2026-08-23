package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	PostgresAddr string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresPort int
	JwtSecret    string
	ServAddr     string
}

func LoadConfigs() (config Config) {
	config.PostgresAddr = getEnv("POSTGRES_ADDRESS")
	config.PostgresUser = getEnv("POSTGRES_USER")
	config.PostgresPass = getEnv("POSTGRES_PASSWORD")
	config.PostgresDB = getEnv("POSTGRES_DB")
	config.PostgresPort, _ = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	config.JwtSecret = getEnv("JWT_SECRET")
	config.ServAddr = getEnv("SERVER_ADDR")
	return config
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
