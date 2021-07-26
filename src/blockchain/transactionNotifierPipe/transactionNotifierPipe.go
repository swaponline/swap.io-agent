package transactionNotifierPipe

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/redisStore"
)

type TransactionNotifierPipe struct {
	input            chan blockchain.Transaction
	Out              chan blockchain.TransactionPipeData
	subscribersStore redisStore.ISubscribersStore
	stop             chan bool
}

type TransactionNotifierPipeConfig struct {
	Input            chan blockchain.Transaction
	SubscribersStore redisStore.ISubscribersStore
}

func InitializeTransactionNotifierPipe(
	config TransactionNotifierPipeConfig,
) *TransactionNotifierPipe {
	return &TransactionNotifierPipe{
		input: config.Input,
		Out: make(chan blockchain.TransactionPipeData),
		subscribersStore: config.SubscribersStore,
	}
}
func (tnp *TransactionNotifierPipe) Start() {
	exit := false
	for !exit {
		select {
			case transaction := <- tnp.input: {
				subscribers := tnp.subscribersStore.GetSubscribersFromAddresses(
					transaction.AllSpendAddresses,
				)
				if len(subscribers) > 0 {
					tnp.Out <- blockchain.TransactionPipeData{
						Subscribers: subscribers,
						Transaction: transaction,
					}
				}
			}
			case <- tnp.stop: {
				exit = true
			}
		}
	}
}
func (tnp *TransactionNotifierPipe) Stop() error {
	tnp.stop <- true
	return nil
}
func (tnp *TransactionNotifierPipe) Status() error {
	return nil
}