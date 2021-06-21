package runApp

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadConfig() error {
	mode := os.Args[0]
	if mode == "production" {
		err := godotenv.Load(".env.production")
		return err
	} else {
		err := godotenv.Load(".env.development")
		return err
	}
}