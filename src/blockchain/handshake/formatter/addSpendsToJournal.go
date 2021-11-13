package transactionFormatter

import (
	"strconv"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
	"swap.io-agent/src/blockchain/journal"
)

func AddSpendsToJournal(
	tx *nodeApi.Transaction,
	journal *journal.Journal,
	minerAddress string,
) {
	feeSize := tx.Fee
	for _, input := range tx.Inputs {
		if input.Coin.Value == 0 {
			continue
		}
		if feeSize > 0 {
			if feeSize >= input.Coin.Value {
				journal.Add(HSD, blockchain.Spend{
					Wallet: input.Coin.Address,
					Value:  strconv.Itoa(-input.Coin.Value),
					Label:  blockchain.SPEND_LABEL_FEE,
				})
				journal.Add(HSD, blockchain.Spend{
					Wallet: minerAddress,
					Value:  strconv.Itoa(input.Coin.Value),
					Label:  blockchain.SPEND_LABEL_FEE,
				})
				feeSize -= input.Coin.Value
			} else {
				journal.Add(HSD, blockchain.Spend{
					Wallet: input.Coin.Address,
					Value:  strconv.Itoa(-feeSize),
					Label:  blockchain.SPEND_LABEL_FEE,
				})
				journal.Add(HSD, blockchain.Spend{
					Wallet: minerAddress,
					Value:  strconv.Itoa(feeSize),
					Label:  blockchain.SPEND_LABEL_FEE,
				})
				journal.Add(HSD, blockchain.Spend{
					Wallet: input.Coin.Address,
					Value:  strconv.Itoa(-(input.Coin.Value - feeSize)),
					Label:  blockchain.SPEND_LABEL_TRANSFER,
				})
				feeSize = 0
			}
		} else {
			journal.Add(HSD, blockchain.Spend{
				Wallet: input.Coin.Address,
				Value:  strconv.Itoa(-input.Coin.Value),
				Label:  blockchain.SPEND_LABEL_TRANSFER,
			})
		}
	}
	for _, output := range tx.Outputs {
		if output.Value > 0 {
			journal.Add(HSD, blockchain.Spend{
				Wallet: output.Address,
				Value:  strconv.Itoa(output.Value),
				Label:  blockchain.SPEND_LABEL_TRANSFER,
			})
		}
	}
}
