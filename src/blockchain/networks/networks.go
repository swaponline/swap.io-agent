package networks

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum"
)

type Networks map[string]blockchain.IBlockchinApi

func InitializeNetworks() *Networks {
	network := Networks{}
	network["ethereum"] = ethereum.InitializeApi()

	return &network
}

func (*Networks) Start() {}
func (*Networks) Stop() error {
	return nil
}
func (*Networks) Status() error {
	return nil
}
