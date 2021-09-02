package transactionFormatter

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
	"swap.io-agent/src/blockchain/journal"
)

const ETH = "ETH"

type TransactionFormatter struct {
	api nodeApi.IGeth
}
type TransactionFormatterConfig struct {
	Api nodeApi.IGeth
}

func InitializeTransactionFormatter(
	config TransactionFormatterConfig,
) *TransactionFormatter {
	return &TransactionFormatter{
		api: config.Api,
	}
}

func (tf *TransactionFormatter) FormatTransactionFromHash(
	hash string,
) (*blockchain.Transaction, error) {
	transaction, err := tf.api.GetTransactionByHash(hash)
	if err != blockchain.ApiRequestSuccess {
		return nil, fmt.Errorf(
			"not get transaction by hash %v", hash,
		)
	}

	transactionBlockIndex, errConv := strconv.ParseInt(transaction.BlockNumber, 0, 64)
	if errConv != nil {
		return nil, errConv
	}
	blockTransaction, err := tf.api.GetBlockByIndex(int(transactionBlockIndex))
	if err != blockchain.ApiRequestSuccess {
		return nil, fmt.Errorf(
			"not get transaction block by index %v", err,
		)
	}

	log.Println("OK")

	return tf.FormatTransaction(transaction, blockTransaction)
}
func (tf *TransactionFormatter) FormatTransaction(
	blockTransaction *nodeApi.BlockTransaction,
	block *nodeApi.Block,
) (*blockchain.Transaction, error) {
	transactionLogs, errReq := tf.api.GetTransactionLogs(
		blockTransaction.Hash,
	)
	if errReq != blockchain.ApiRequestSuccess {
		return nil, fmt.Errorf(
			"not get transactionLogs error - %v", errReq,
		)
	}
	//internalTransaction, errReq := tf.api.GetInternalTransaction(
	//	blockTransaction.Hash,
	//)
	//if errReq != api.RequestSuccess {
	//	return nil, fmt.Errorf(
	//		"not get interanlTransaction error - %v", errReq,
	//	)
	//}

	transactionGasUsed, ok := new(big.Int).SetString(
		transactionLogs.GasUsed, 0,
	)
	if !ok {
		if bytes, err := json.Marshal(transactionLogs); err != nil {
			log.Println(string(bytes))
		}
		return nil, fmt.Errorf(
			"transactionLogs.Result.GasUsed(%v) not parsed %v",
			transactionLogs.GasUsed,
			ok,
		)
	}
	transactionGasPrice, ok := new(big.Int).SetString(
		blockTransaction.GasPrice, 0,
	)
	if !ok {
		return nil, fmt.Errorf(
			"blockTransaction.GasPrice(%v) not parsed",
			blockTransaction.GasPrice,
		)
	}

	transactionFee := big.NewInt(0).Mul(
		transactionGasUsed, transactionGasPrice,
	).Text(16)

	transactionJournal := journal.New("ethereum")
	if len(blockTransaction.From) > 0 {
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: blockTransaction.From,
			Value:  `-` + blockTransaction.Value,
		})
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: blockTransaction.From,
			Value:  `-` + transactionFee,
			Label:  "Transaction fee",
		})
	}
	if len(block.Miner) > 0 {
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: block.Miner,
			Value:  transactionFee,
			Label:  "Transaction fee",
		})
	}
	if len(blockTransaction.To) > 0 {
		transactionJournal.Add(ETH, blockchain.Spend{
			Wallet: blockTransaction.To,
			Value:  blockTransaction.Value,
		})
	}

	//AddSpendsFromInternalTxCallsToJournal(
	//	internalTransaction.Calls,
	//	transactionJournal,
	//)
	AddSpendsFromLogsToJournal(
		transactionLogs.Logs,
		transactionJournal,
	)

	transaction := blockchain.Transaction{
		Hash:              blockTransaction.Hash,
		From:              blockTransaction.From,
		To:                blockTransaction.To,
		Gas:               blockTransaction.Gas,
		GasPrice:          blockTransaction.GasPrice,
		GasUsed:           transactionLogs.GasUsed,
		Value:             blockTransaction.Value,
		Timestamp:         block.Timestamp,
		TransactionIndex:  blockTransaction.TransactionIndex,
		BlockHash:         blockTransaction.BlockHash,
		BlockNumber:       blockTransaction.BlockNumber,
		BlockMiner:        block.Miner,
		Nonce:             blockTransaction.Nonce,
		AllSpendAddresses: transactionJournal.GetSpendsAddress(),
		Journal:           transactionJournal.GetSpends(),
	}

	return &transaction, nil
}

func (*TransactionFormatter) Start() {}
func (*TransactionFormatter) Stop() error {
	return nil
}
func (*TransactionFormatter) Status() error {
	return nil
}
