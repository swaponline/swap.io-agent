package transactionFormater

import (
	"errors"
	"fmt"
	"strconv"
	"swap.io-agent/src/blockchainHandlers"
	"swap.io-agent/src/blockchainHandlers/ethereum/api/ethercsan"
)

func FormatTransaction(
	apiKey string,
	blockTransaction *ethercsan.BlockTransaction,
	block ethercsan.Block,
) (*blockchainHandlers.Transaction, error) {
	transactionLogs, errReq := ethercsan.GetTransactionLogs(
		apiKey,
		blockTransaction.Hash,
	)
	if errReq != ethercsan.RequestSuccess {
		return nil, errors.New(fmt.Sprintf(
			"not get transactionLogs error - %v", errReq,
		))
	}
	blockTransactionNonce, err := strconv.ParseInt(
		block.Nonce,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}
	blockTransactionTimestamp, err := strconv.ParseInt(
		block.Timestamp,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}
	blockTransactionIndex, err := strconv.ParseInt(
		blockTransaction.Value,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}
	blockTransactionValue, err := strconv.ParseInt(
		blockTransaction.Value,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}
	blockTransactionGas, err := strconv.ParseInt(
		blockTransaction.Gas,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}
	blockTransactionGasPrice, err := strconv.ParseInt(
		blockTransaction.GasPrice,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}
	blockTransactionGasUsed, err := strconv.ParseInt(
		transactionLogs.Result.GasUsed,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}
	blockTransactionBlock, err := strconv.ParseInt(
		blockTransaction.BlockNumber,
		16,
		64,
	)
	if err != nil {
		return nil, err
	}

	transaction := blockchainHandlers.Transaction{
		Hash: blockTransaction.Hash,
		From: blockTransaction.From,
		To:   blockTransaction.To,
		Gas: int(blockTransactionGas),
		GasPrice: int(blockTransactionGasPrice),
		GasUsed: int(blockTransactionGasUsed),
		Value: int(blockTransactionValue),
		Timestamp: int(blockTransactionTimestamp),
		TransactionIndex: int(blockTransactionIndex),
		BlockHash: block.Hash,
		BlockNumber: int(blockTransactionBlock),
		BlockMiner: block.Miner,
		Nonce: int(blockTransactionNonce),
		AllSpendAddresses: ethercsan.GetAllSpendAddressFromLogs(
			transactionLogs.Result.Logs,
			blockTransaction,
			block.Miner,
		),
	}

	return &transaction, nil
}
