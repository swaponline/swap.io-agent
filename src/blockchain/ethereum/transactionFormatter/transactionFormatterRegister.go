package transactionFormatter

import (
	"log"
	"swap.io-agent/src/blockchain/ethereum/api/geth"
	"swap.io-agent/src/serviceRegistry"
)


func Register(reg *serviceRegistry.ServiceRegistry) {
	var api *geth.Geth
	err := reg.FetchService(&api)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeTransactionFormatter(TransactionFormatterConfig{
			Api: api,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}