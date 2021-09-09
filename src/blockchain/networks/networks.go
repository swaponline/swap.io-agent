package networks

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake"
)

type Networks map[string]blockchain.IBlockchinApi

func InitializeNetworks() *Networks {
	network := Networks{}
	//network["ethereum"] = ethereum.InitializeApi()
	network["handshake"] = handshake.InitializeApi()

	return &network
}

func (*Networks) Start() {}
func (*Networks) Stop() error {
	return nil
}
func (*Networks) Status() error {
	return nil
}
