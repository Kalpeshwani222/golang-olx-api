package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Env       string
	DbUrl     string
	JWTSecret string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		panic("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET is required")
	}

	return Config{
		Port:      port,
		Env:       env,
		DbUrl:     dbUrl,
		JWTSecret: jwtSecret,
	}
}
