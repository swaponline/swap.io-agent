package env

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var PORT int
var BLOCKCHAIN string
var BLOCKCHAIN_DEFAULT_SCANNED_BLOCK int
var BLOCKCHAIN_REQUESTS_LIMIT int
var SECRET_TOKEN string
var REDIS_ADDR string
var REDIS_PASSWORD string
var REDIS_DB int
var ETHERSCAN_API_KEY string

func InitializeConfig() error {
	mode := os.Args[1]
	log.Printf("Starting width mode - %v", mode)

	if mode == "production" {
		err := godotenv.Load(".env.production")
		if err != nil {
			return err
		}
	} else {
		err := godotenv.Load(".env.development")
		if err != nil {
			return err
		}
	}

	PortInt, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return errors.New("SET PORT IN ENV")
	}
	PORT = PortInt

	BLOCKCHAIN = os.Getenv("BLOCKCHAIN")
	if len(BLOCKCHAIN) == 0 {
		return errors.New("SET BLOCKCHAIN IN ENV")
	}

	blockchainDefaultScannedBlockInt, err := strconv.Atoi(
		os.Getenv("BLOCKCHAIN_DEFAULT_SCANNED_BLOCK"),
	)
	if err != nil {
		return errors.New("SET BLOCKCHAIN_DEFAULT_SCANNED_BLOCK IN ENV")
	}
	BLOCKCHAIN_DEFAULT_SCANNED_BLOCK = blockchainDefaultScannedBlockInt

	blockchainRequestsLimitInt, err := strconv.Atoi(
		os.Getenv("BLOCKCHAIN_REQUESTS_LIMIT"),
	)
	if err != nil {
		return errors.New("SET BLOCKCHAIN_REQUESTS_LIMIT IN ENV")
	}
	BLOCKCHAIN_REQUESTS_LIMIT = blockchainRequestsLimitInt

	SECRET_TOKEN = os.Getenv("SECRET_TOKEN")
	if len(SECRET_TOKEN) == 0 {
		return errors.New("SET SECRET_TOKEN IN ENV")
	}

	REDIS_ADDR = os.Getenv("REDIS_ADDR")
	if len(REDIS_ADDR) == 0 {
		return errors.New("SET REDIS_ADDR IN ENV")
	}

	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	if len(REDIS_PASSWORD) == 0 {
		return errors.New("SET REDIS_PASSWORD IN ENV")
	}

	redisDbInt, err := strconv.Atoi(
		os.Getenv("REDIS_DB"),
	)
	if err != nil {
		return errors.New("SET REDIS_DB IN ENV")
	}
	REDIS_DB = redisDbInt

	ETHERSCAN_API_KEY = os.Getenv("ETHERSCAN_API_KEY")
	if len(ETHERSCAN_API_KEY) == 0 {
		return errors.New("SET ETHERSCAN_API_KEY IN ENV")
	}

	return nil
}