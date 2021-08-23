package levelDbStore

type ITransactionsStore interface {
	GetLastTransactionBlock() int
	GetAddressTransactionsHash(
		address string,
		startTime int,
		endTime int,
	) ([]string, error)
	WriteLastIndexedTransactions(
		AddressHashTransactions map[string][]string,
		indexBlock int,
	) error
	GetCursorFromAddress(
		address string,
	) (string, error)
	GetCursorTransactionHashes(
		cursor string,
	) (*CursorTransactionHashes, error)
	GetFirstCursorTransactionHashes(
		address string,
	) (*CursorTransactionHashes,error)
}
