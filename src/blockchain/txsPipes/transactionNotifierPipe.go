package txsPipes

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/subscribersManager"
)

type NewTxsNotifierPipe struct {
	inp                chan *blockchain.Transaction
	Out                chan *blockchain.TransactionPipeData
	subscribersManager *subscribersManager.SubscribesManager
	stop               chan struct{}
}

type NewTxsNotifierPipeConfig struct {
	Input              chan *blockchain.Transaction
	SubscribersManager *subscribersManager.SubscribesManager
}

func InitializeNewTxsNotifierPipe(
	config NewTxsNotifierPipeConfig,
) *NewTxsNotifierPipe {
	return &NewTxsNotifierPipe{
		inp:                config.Input,
		Out:                make(chan *blockchain.TransactionPipeData),
		subscribersManager: config.SubscribersManager,
	}
}
func (ntnp *NewTxsNotifierPipe) Start() {
	FilterExpectedTxsPipe(
		ntnp.subscribersManager,
		ntnp.inp,
		ntnp.Out,
		ntnp.stop,
	)
}
func (ntnp *NewTxsNotifierPipe) Stop() error {
	close(ntnp.stop)
	return nil
}
func (ntnp *NewTxsNotifierPipe) Status() error {
	return nil
}
