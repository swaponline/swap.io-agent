package fullNodeApi

import (
	"log"
	"net/http"
	"os"
	transactionFormatter "swap.io-agent/src/blockchain/handshake/formatter"
)

type FullNodeApi struct {
	baseUrl   string
	apiKey    string
	formatter transactionFormatter.TransactionFormatter
	client    http.Client
}

func InitializeFullNodeApi() *FullNodeApi {
	baseUrl := os.Getenv("HANDSHAKE_BASE_URL")
	if len(baseUrl) == 0 {
		log.Panicln("SET HANDSHAKE_BASE_URL IN ENV")
	}
	apiKey := os.Getenv("HANDSHAKE_API_KEY")
	if len(apiKey) == 0 {
		log.Panicln("SET HANDSHAKE_API_KEY IN ENV")
	}

	return &FullNodeApi{
		baseUrl: baseUrl,
		apiKey:  apiKey,
		client:  http.Client{},
	}
}
