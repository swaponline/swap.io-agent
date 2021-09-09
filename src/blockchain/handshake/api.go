package handshake

import (
	transactionFormatter "swap.io-agent/src/blockchain/handshake/formatter"
	"swap.io-agent/src/blockchain/handshake/nodeApi/fullNodeApi"
)

type Api struct {
	nodeApi   *fullNodeApi.FullNodeApi
	formatter *transactionFormatter.TransactionFormatter
}

func Test() {
	InitializeApi().nodeApi.GetBlockByIndex(0)
}

func InitializeApi() *Api {
	fullNodeApi := fullNodeApi.InitializeFullNodeApi()
	formatterNodeApi := transactionFormatter.InitializeTransactionFormatter(
		transactionFormatter.TransactionFormatterConfig{
			Api: fullNodeApi,
		},
	)

	return &Api{
		nodeApi:   fullNodeApi,
		formatter: formatterNodeApi,
	}
}
