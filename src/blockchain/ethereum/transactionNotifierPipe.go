package ethereum

import "swap.io-agent/src/blockchain"

type subscribersStore interface {
	GetSubscribersFromAddresses(addresses []string) []string
}

type TransactionNotifierPipe struct {
	input chan blockchain.Transaction
	Out   chan blockchain.TransactionPipeData
	subscribersStore subscribersStore
	stop  chan bool
}

type TransactionNotifierPipeConfig struct {
	Input chan blockchain.Transaction
	SubscribersStore subscribersStore
}

func TransactionNotifierPipeInitialize(
	config TransactionNotifierPipeConfig,
) TransactionNotifierPipe {
	return TransactionNotifierPipe{
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