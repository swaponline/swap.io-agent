package levelDbStore

type ITransactionsStore interface {
	GetLastTransactionBlock() int
	GetAddressTransactionsHash(
		address string,
		startTime int,
		endTime int,
	) ([]string, error)
	WriteLastIndexedBlockTransactions(
		indexedTransactions *map[string][]string,
		indexBlock int,
		timestampBlock int,
	) error
}
