package ethereum

import (
	"log"
	"swap.io-agent/src/blockchainHandlers/ethereum/api/ethercsan"
	"time"
)



func (indexer *BlockchainIndexer) runScanner() {
	isSynchronize := false
	for {
		block, err := ethercsan.GetBlockByIndex(
			indexer.apiKey,
			indexer.lastBlock,
		)
		if err == ethercsan.NotExistBlock {
			if !isSynchronize {
				isSynchronize = true
				close(indexer.isSynchronize)
			}

			<-time.After(time.Minute)
			continue
		}
		if err != ethercsan.SuccessReq {
			log.Panicf("not get block by index request etherscan")
		}

		if isSynchronize && len(block.Transactions) > 0 {

		}
	}
}