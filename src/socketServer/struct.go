package socketServer

import "swap.io-agent/src/blockchain"

type SubscribeEventPayload struct {
	address string
	startTime int
	endTime int
}
type SynchroniseAddressData struct {
	Address string `json:"address"`
	Transactions []*blockchain.Transaction `json:"transactions"`
}