package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var PORT int
var BLOCKCHAIN string
var BLOCKCHAIN_DEFAULT_SCANNED_BLOCK int
var BLOCKCHAIN_REQUESTS_LIMIT int
var SECRET_TOKEN string
var KAFKA_ADDR string
var REDIS_ADDR string
var REDIS_PASSWORD string
var REDIS_DB int
var ETHERSCAN_API_KEY string

func InitializeConfig() error {
	godotenv.Load()
	PortInt, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Panicln("SET PORT IN ENV")
	}
	PORT = PortInt

	BLOCKCHAIN = os.Getenv("BLOCKCHAIN")
	if len(BLOCKCHAIN) == 0 {
		log.Panicln("SET BLOCKCHAIN IN ENV")
	}

	blockchainDefaultScannedBlockInt, err := strconv.Atoi(
		os.Getenv("BLOCKCHAIN_DEFAULT_SCANNED_BLOCK"),
	)
	if err != nil {
		log.Panicln("SET BLOCKCHAIN_DEFAULT_SCANNED_BLOCK IN ENV")
	}
	BLOCKCHAIN_DEFAULT_SCANNED_BLOCK = blockchainDefaultScannedBlockInt

	blockchainRequestsLimitInt, err := strconv.Atoi(
		os.Getenv("BLOCKCHAIN_REQUESTS_LIMIT"),
	)
	if err != nil {
		log.Panicln("SET BLOCKCHAIN_REQUESTS_LIMIT IN ENV")
	}
	BLOCKCHAIN_REQUESTS_LIMIT = blockchainRequestsLimitInt

	SECRET_TOKEN = os.Getenv("SECRET_TOKEN")
	if len(SECRET_TOKEN) == 0 {
		log.Panicln("SET SECRET_TOKEN IN ENV")
	}

	KAFKA_ADDR = os.Getenv("KAFKA_ADDR")
	if len(KAFKA_ADDR) == 0 {
		log.Panicln("SET KAFKA_ADDR IN ENV")
	}

	REDIS_ADDR = os.Getenv("REDIS_ADDR")
	if len(REDIS_ADDR) == 0 {
		log.Panicln("SET REDIS_ADDR IN ENV")
	}

	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	if len(REDIS_PASSWORD) == 0 {
		log.Panicln("SET REDIS_PASSWORD IN ENV")
	}

	redisDbInt, err := strconv.Atoi(
		os.Getenv("REDIS_DB"),
	)
	if err != nil {
		log.Panicln("SET REDIS_DB IN ENV")
	}
	REDIS_DB = redisDbInt

	ETHERSCAN_API_KEY = os.Getenv("ETHERSCAN_API_KEY")
	if len(ETHERSCAN_API_KEY) == 0 {
		log.Panicln("SET ETHERSCAN_API_KEY IN ENV")
	}

	return nil
}
