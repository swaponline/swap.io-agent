package blockchain

import "swap.io-agent/src/levelDbStore"

type Synchronizer struct {
	store levelDbStore.ITransactionsStore
}
type SynchronizerConfig struct {
	Store levelDbStore.ITransactionsStore
}

func InitialiseSynchronizer(config SynchronizerConfig) *Synchronizer {
	return &Synchronizer{
		store: config.Store,
	}
}

func (s *Synchronizer) SynchronizeAddress(
	address string,
	startTime int,
	endTime int,
)([]Transaction, error) {
	//transactionsHash, err := s.store.GetAddressTransactionsHash(
	//	address,
	//	startTime,
	//	endTime,
	//)
	//if err != nil {
	//	return nil, err
	//}

	return nil,nil
}

func Start() {}
func Stop() error {
	return nil
}
func Status() error {
	return nil
}