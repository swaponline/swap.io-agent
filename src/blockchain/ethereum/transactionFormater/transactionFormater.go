package transactionFormater

import (
	"errors"
	"fmt"
	"math/big"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/api/ethercsan"
	journal "swap.io-agent/src/blockchain/journal"
)

const ETH = "ETH"

func FormatTransaction(
	apiKey string,
	blockTransaction *ethercsan.BlockTransaction,
	block *ethercsan.Block,
) (*blockchain.Transaction, error) {
	transactionLogs, errReq := ethercsan.GetTransactionLogs(
		apiKey,
		blockTransaction.Hash,
	)
	if errReq != ethercsan.RequestSuccess {
		return nil, errors.New(fmt.Sprintf(
			"not get transactionLogs error - %v", errReq,
		))
	}

	transactionGasUsed, ok  := new(big.Int).SetString(
		transactionLogs.Result.GasUsed, 0,
	)
	if !ok {
		return nil, errors.New("transactionLogs.Result.GasUsed not parsed")
	}
	transactionGasPrice, ok := new(big.Int).SetString(
		blockTransaction.GasPrice, 0,
	)
	if !ok {
		return nil, errors.New("blockTransaction.GasPrice not parsed")
	}

	transactionFee := big.NewInt(0).Mul(
		transactionGasUsed, transactionGasPrice,
	).Text(16)

	transactionJournal := journal.New("ethereum")
	transactionJournal.Add(ETH, blockchain.Spend{
		Wallet: blockTransaction.From,
		Value: `-`+blockTransaction.Value,
	})
	transactionJournal.Add(ETH, blockchain.Spend{
		Wallet: blockTransaction.From,
		Value: `-`+transactionFee,
	})
	transactionJournal.Add(ETH, blockchain.Spend{
		Wallet: block.Miner,
		Value: transactionFee,
	})
	transactionJournal.Add(ETH, blockchain.Spend{
		Wallet: blockTransaction.To,
		Value: blockTransaction.Value,
	})

	err := ethercsan.AddSpendsFromLogsToJournal(
		transactionLogs.Result.Logs,
		transactionJournal,
	)
	if err != nil {
		return nil, err
	}

	transaction := blockchain.Transaction{
		Hash: blockTransaction.Hash,
		From: blockTransaction.From,
		To:   blockTransaction.To,
		Gas:  blockTransaction.Gas,
		GasPrice: blockTransaction.GasPrice,
		GasUsed: transactionLogs.Result.GasUsed,
		Value: blockTransaction.Value,
		Timestamp: block.Timestamp,
		TransactionIndex: blockTransaction.TransactionIndex,
		BlockHash: blockTransaction.BlockHash,
		BlockNumber: blockTransaction.BlockNumber,
		BlockMiner: block.Miner,
		Nonce: blockTransaction.Nonce,
		AllSpendAddresses: transactionJournal.GetSpendsAddress(),
		Journal: transactionJournal.GetSpends(),
	}

	return &transaction, nil
}
