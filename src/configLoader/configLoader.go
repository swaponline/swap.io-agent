package configLoader

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitializeConfig() error {
	mode := os.Args[1]
	log.Printf("Starting width mode - %v", mode)

	if mode == "production" {
		err := godotenv.Load(".env.production")
		return err
	} else {
		err := godotenv.Load(".env.development")
		return err
	}
}