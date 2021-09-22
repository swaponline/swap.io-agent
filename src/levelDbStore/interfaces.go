package levelDbStore

type ITransactionsStore interface {
	GetLastTransactionBlock() int
	GetCursorFromAddress(
		address string,
	) (string, error)
	GetFirstCursorTransactionHashes(
		address string,
	) (*CursorTransactionHashes, error)
	GetCursorTransactionHashes(
		cursor string,
	) (*CursorTransactionHashes, error)
	WriteLastIndexedTransactions(
		AddressHashTransactions map[string][]string,
		indexBlock int,
	) error
	Flush() error
}
