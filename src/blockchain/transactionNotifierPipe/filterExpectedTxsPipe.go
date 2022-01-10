package transactionNotifierPipe

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/subscribersManager"
)

func FilterExpectedTxsPipe(
	subscribersManager *subscribersManager.SubscribesManager,
	inp chan *blockchain.Transaction,
	out chan *blockchain.TransactionPipeData,
	stop chan struct{},
) {
	for {
		select {
		case transaction := <-inp:
			{
				subscribers := subscribersManager.GetSubscribersFromAddresses(
					transaction.AllSpendAddresses,
				)
				if len(subscribers) > 0 {
					out <- &blockchain.TransactionPipeData{
						Subscribers: subscribers,
						Transaction: transaction,
					}
				}
			}

		case <-stop:
			return
		}
	}
}
