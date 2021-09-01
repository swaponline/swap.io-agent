package ethercsan

import (
	"log"
	"os"
)

type Etherscan struct {
	apiKey string
	baseUrl string
}

func InitializeEthercsan() *Etherscan {
	etherscanApiKey := os.Getenv("ETHERSCAN_API_KEY")
	if len(etherscanApiKey) == 0 {
		log.Panicln("SET ETHERSCAN_API_KEY IN ENV")
	}

	etherscanBaseUrl := os.Getenv("ETHERSCAN_BASE_URL")
	if len(etherscanBaseUrl) == 0 {
		log.Panicln("SET ETHERSCAN_BASE_URL IN ENV")
	}

	return &Etherscan{
		apiKey: etherscanApiKey,
		baseUrl: etherscanBaseUrl,
	}
}

func(*Etherscan) Start() {}
func(*Etherscan) Stop() error {
	return nil
}
func(*Etherscan) Status() error {
	return nil
}