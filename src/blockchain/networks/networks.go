package networks

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum"
)

type Network map[string]blockchain.IBlockchinApi

func IndexerRegister() *Network {
	network := Network{}
	network["ethereum"] = ethereum.InitializeApi()

	return &network
}

func (*Network) Start() {}
func (*Network) Stop() error {
	return nil
}
func (*Network) Status() error {
	return nil
}
