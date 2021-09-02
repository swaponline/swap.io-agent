package ethereum

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
	"swap.io-agent/src/blockchain/ethereum/nodeApi/geth"
	"swap.io-agent/src/blockchain/ethereum/transactionFormatter"
)

type Api struct {
	nodeApi   nodeApi.IGeth
	formatter blockchain.IFormatter
}

func InitializeApi() *Api {
	geth := geth.InitializeGeth()
	gethFormatter := transactionFormatter.InitializeTransactionFormatter(
		transactionFormatter.TransactionFormatterConfig{
			Api: geth,
		},
	)

	return &Api{
		nodeApi:   geth,
		formatter: gethFormatter,
	}
}
