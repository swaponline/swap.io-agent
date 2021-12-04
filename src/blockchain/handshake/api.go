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
	//InitializeApi().nodeApi.GetTransactionByHash("8f99b0037eb07812737aaa1005af85fc4429e20a65f66bf15d148be02abca587")
}

func InitializeApi() *Api {
	fullNodeApiInstance := fullNodeApi.InitializeFullNodeApi()
	formatterNodeApi := transactionFormatter.InitializeTransactionFormatter(
		transactionFormatter.TransactionFormatterConfig{
			Api: fullNodeApiInstance,
		},
	)

	return &Api{
		nodeApi:   fullNodeApiInstance,
		formatter: formatterNodeApi,
	}
}
