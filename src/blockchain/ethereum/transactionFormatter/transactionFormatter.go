package transactionFormatter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum"
	"swap.io-agent/src/blockchain/ethereum/api"
	journal "swap.io-agent/src/blockchain/journal"
)

const ETH = "ETH"

type TransactionFormatter struct {
	api ethereum.IGeth
}
type TransactionFormatterConfig struct {
	Api ethereum.IGeth
}

func InitializeTransactionFormatter(config TransactionFormatterConfig)*TransactionFormatter {
	return &TransactionFormatter{
		api: config.Api,
	}
}

func (tf *TransactionFormatter) FormatTransactionFromHash(
	hash string,
) (*blockchain.Transaction, error) {
	transaction, err := tf.api.GetTransactionByHash(hash)
	if err != api.RequestSuccess {
		return nil, errors.New(
			fmt.Sprintf("not get transaction by hash %v", hash),
		)
	}

	transactionBlockIndex, errConv := strconv.Atoi(transaction.BlockNumber)
	if errConv != nil {
		return nil, errConv
	}
	blockTransaction, err := tf.api.GetBlockByIndex(transactionBlockIndex)
	if err != api.RequestSuccess {
		return nil, errors.New(
			fmt.Sprintf("not get transaction block by index %v", err),
		)
	}

	return tf.FormatTransaction(transaction, blockTransaction)
}
func (tf *TransactionFormatter) FormatTransaction(
	blockTransaction *api.BlockTransaction,
	block *api.Block,
) (*blockchain.Transaction, error) {
	transactionLogs, errReq := tf.api.GetTransactionLogs(
		blockTransaction.Hash,
	)
	if errReq != api.RequestSuccess {
		return nil, errors.New(fmt.Sprintf(
			"not get transactionLogs error - %v", errReq,
		))
	}

	transactionGasUsed, ok  := new(big.Int).SetString(
		transactionLogs.GasUsed, 0,
	)
	if !ok {
		if bytes, err := json.Marshal(transactionLogs); err != nil {
			log.Println(string(bytes))
		}
		return nil, errors.New(fmt.Sprintf(
			"transactionLogs.Result.GasUsed(%v) not parsed %v",
			transactionLogs.GasUsed,
			ok,
		))
	}
	transactionGasPrice, ok := new(big.Int).SetString(
		blockTransaction.GasPrice, 0,
	)
	if !ok {
		return nil, errors.New(
			fmt.Sprintf(
				"blockTransaction.GasPrice(%v) not parsed",
				blockTransaction.GasPrice,
			),
		)
	}

	transactionFee := big.NewInt(0).Mul(
		transactionGasUsed, transactionGasPrice,
	).Text(16)

	transactionJournal := journal.New("ethereum")
	//if len(blockTransaction.From) == 0 {
	//	log.Panicln(blockTransaction.Hash, blockTransaction.From, "FROM")
	//}
	//if len(block.Miner) == 0 {
	//	log.Panicln(blockTransaction.Hash, block.Miner, "MINER")
	//}
	//if len(blockTransaction.To) == 0 {
	//	log.Panicln(blockTransaction.Hash, block.Number, blockTransaction.To, "TO")
	//}
	if len(blockTransaction.From) > 0 {
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: blockTransaction.From,
			Value: `-`+blockTransaction.Value,
		})
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: blockTransaction.From,
			Value: `-`+transactionFee,
			Label: "Transaction fee",
		})
	}
	if len(block.Miner) > 0 {
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: block.Miner,
			Value: transactionFee,
			Label: "Transaction fee",
		})
	}
	if len(blockTransaction.To) > 0 {
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: blockTransaction.To,
			Value: blockTransaction.Value,
		})
	}

	err := AddSpendsFromLogsToJournal(
		transactionLogs.Logs,
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
		GasUsed: transactionLogs.GasUsed,
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

func (_ *TransactionFormatter) Start() {}
func (_ *TransactionFormatter) Stop() error {
	return nil
}
func (_ *TransactionFormatter) Status() error {
	return nil
}