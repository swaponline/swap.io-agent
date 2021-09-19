package socketServer

import "swap.io-agent/src/blockchain"

type SubscribeEventPayload struct {
	Address string `json:"address"`
}
type SynchroniseAddressData struct {
	Address      string                    `json:"address"`
	Transactions []*blockchain.Transaction `json:"transactions"`
}
