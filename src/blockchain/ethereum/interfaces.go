package ethereum

import "swap.io-agent/src/blockchain/ethereum/nodeApi"

type IGeth interface {
	GetBlockByIndex(
		index int,
	) (*nodeApi.Block, int)
	GetTransactionByHash(
		hash string,
	) (*nodeApi.BlockTransaction, int)
	GetTransactionLogs(
		hash string,
	) (*nodeApi.TransactionLogs, int)
	GetInternalTransaction(
		hash string,
	) (*nodeApi.InteranlTransaction, int)
	GetBlockTraceByIndex(
		index string,
	) (interface{}, int)
}
