package transactionNotifierPipe

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/subscribersManager"
)

type TransactionNotifierPipe struct {
	inp                chan *blockchain.Transaction
	Out                chan *blockchain.TransactionPipeData
	subscribersManager *subscribersManager.SubscribesManager
	stop               chan struct{}
}

type TransactionNotifierPipeConfig struct {
	Input              chan *blockchain.Transaction
	SubscribersManager *subscribersManager.SubscribesManager
}

func InitializeTransactionNotifierPipe(
	config TransactionNotifierPipeConfig,
) *TransactionNotifierPipe {
	return &TransactionNotifierPipe{
		inp:                config.Input,
		Out:                make(chan *blockchain.TransactionPipeData),
		subscribersManager: config.SubscribersManager,
	}
}
func (tnp *TransactionNotifierPipe) Start() {
	FilterExpectedTxsPipe(tnp.subscribersManager, tnp.inp, tnp.Out, tnp.stop)
}
func (tnp *TransactionNotifierPipe) Stop() error {
	close(tnp.stop)
	return nil
}
func (tnp *TransactionNotifierPipe) Status() error {
	return nil
}
