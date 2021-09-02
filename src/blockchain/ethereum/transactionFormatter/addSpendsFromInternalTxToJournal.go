package transactionFormatter

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
	"swap.io-agent/src/blockchain/journal"
)

func AddSpendsFromInternalTxCallsToJournal(
	internalTxCalls []nodeApi.InteranlTransactionCall,
	journal *journal.Journal,
) {
	for _, call := range internalTxCalls {
		// 0x 0x0 - skip | 0x1 add
		if IsNotEmptyVal(call.Value) {
			if len(call.From) > 0 {
				journal.Add(call.From, blockchain.Spend{
					Wallet: call.From,
					Value:  "-" + call.Value,
				})
			}
			if len(call.To) > 0 {
				journal.Add(call.To, blockchain.Spend{
					Wallet: call.To,
					Value:  call.Value,
				})
			}
		}
	}
}
