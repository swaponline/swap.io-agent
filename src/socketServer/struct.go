package socketServer

import "swap.io-agent/src/blockchain"

type SubscribeEventPayload struct {
	Address string `json:"address"`
	StartTime int  `json:"start_time"`
}
type SynchroniseAddressData struct {
	Address string `json:"address"`
	Transactions []*blockchain.Transaction `json:"transactions"`
}