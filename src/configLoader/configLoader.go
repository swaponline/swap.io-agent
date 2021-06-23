package configLoader

import (
	"github.com/joho/godotenv"
	"os"
)

func InitializeConfig() error {
	mode := os.Args[0]
	if mode == "production" {
		err := godotenv.Load(".env.production")
		return err
	} else {
		err := godotenv.Load(".env.development")
		return err
	}
}