package envs

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {

	envPath, err := filepath.Abs(filepath.Join("..", ".env"))
	if err != nil {
		log.Fatal("Error determining .env file path")
	}

	err = godotenv.Load(envPath)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
