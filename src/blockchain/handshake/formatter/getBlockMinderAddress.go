package transactionFormatter

import (
	"log"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
)

func GetBlockMinderAddress(
	block *nodeApi.Block,
) (
	address string,
) {
	for _, tx := range block.Txs {
		if tx.Fee == 0 {
			address = tx.Outputs[0].Address
			return address
		}
	}
	log.Panicln(block.Hash, "not found miner")
	return ""
}
