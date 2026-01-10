package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
}

var Cfg Config
var cfgOnce sync.Once

func GetConfig() Config {
	cfgOnce.Do(func() {
		instance, err := loadConfig()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		Cfg = instance
	})
	return Cfg
}

func loadConfig() (Config, error) {
	err := godotenv.Load(filepath.Join("..", ".env"))
	return Config{
		DatabaseURL: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}, err
}
