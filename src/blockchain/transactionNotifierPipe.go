package blockchain

import "swap.io-agent/src/redisStore"

type TransactionNotifierPipe struct {
	input            chan Transaction
	Out              chan TransactionPipeData
	subscribersStore redisStore.SubscribersStore
	stop             chan bool
}

type TransactionNotifierPipeConfig struct {
	Input            chan Transaction
	SubscribersStore redisStore.SubscribersStore
}

func InitializeTransactionNotifierPipe(
	config TransactionNotifierPipeConfig,
) TransactionNotifierPipe {
	return TransactionNotifierPipe{
		input: config.Input,
		Out: make(chan TransactionPipeData),
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
					tnp.Out <- TransactionPipeData{
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