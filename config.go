package main

import (
	"os"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	JWTSecret  string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Port:       getEnv("PORT", ":8000"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBAddress:  getEnv("DB_ADDRESS", "localhost:8000"),
		DBName:     getEnv("DB_NAME", "Project_Manager"),
		JWTSecret:  getEnv("JWT_SECRET", "randomsecretkey"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
