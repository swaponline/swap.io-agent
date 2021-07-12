package ethereum

import "swap.io-agent/src/blockchain"

type SubscribersStore interface {
	GetSubscribersFromAddresses(addresses []string) []string
}

type TransactionNotifierPipe struct {
	Input chan blockchain.Transaction
	Out   chan blockchain.TransactionPipeData
	SubscribersStore SubscribersStore
	stop  chan bool
}

type TransactionNotifierPipeConfig struct {
	Input chan blockchain.Transaction
	SubscribersStore SubscribersStore
}

func TransactionNotifierPipeInitialize(
	config TransactionNotifierPipeConfig,
) TransactionNotifierPipe {
	return TransactionNotifierPipe{
		Input: config.Input,
		Out: make(chan blockchain.TransactionPipeData),
		SubscribersStore: config.SubscribersStore,
	}
}
func (tnp *TransactionNotifierPipe) Start() {
	exit := false
	for !exit {
		select {
			case transaction := <- tnp.Input: {
				subscribers := tnp.SubscribersStore.GetSubscribersFromAddresses(
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