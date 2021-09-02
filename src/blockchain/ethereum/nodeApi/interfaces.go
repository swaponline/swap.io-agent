package nodeApi

type IGeth interface {
	GetBlockByIndex(
		index int,
	) (*Block, int)
	GetTransactionByHash(
		hash string,
	) (*BlockTransaction, int)
	GetTransactionLogs(
		hash string,
	) (*TransactionLogs, int)
	GetInternalTransaction(
		hash string,
	) (*InteranlTransaction, int)
	GetBlockTraceByIndex(
		index string,
	) (interface{}, int)
}
