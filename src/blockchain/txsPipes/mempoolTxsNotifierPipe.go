package txsPipes

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/subscribersManager"
)

type MempoolTxsNotifierPipe struct {
	inp                chan *blockchain.Transaction
	Out                chan *blockchain.TransactionPipeData
	subscribersManager *subscribersManager.SubscribesManager
	stop               chan struct{}
}

type MempoolTxsNotifierPipeConfig struct {
	Input              chan *blockchain.Transaction
	SubscribersManager *subscribersManager.SubscribesManager
}

func InitializeMempoolTxsNotifierPipe(
	config MempoolTxsNotifierPipeConfig,
) *MempoolTxsNotifierPipe {
	return &MempoolTxsNotifierPipe{
		inp:                config.Input,
		Out:                make(chan *blockchain.TransactionPipeData),
		subscribersManager: config.SubscribersManager,
	}
}
func (mtnp *MempoolTxsNotifierPipe) Start() {
	FilterExpectedTxsPipe(
		mtnp.subscribersManager,
		mtnp.inp,
		mtnp.Out,
		mtnp.stop,
	)
}
func (mtnp *MempoolTxsNotifierPipe) Stop() error {
	close(mtnp.stop)
	return nil
}
func (mtnp *MempoolTxsNotifierPipe) Status() error {
	return nil
}
