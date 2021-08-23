package levelDbStore

type CursorTransactionHashes struct {
	Cursor     string
	NextCursor string
	Hashes     []string
}