package blockchain

const (
	ApiRequestSuccess = iota
	ApiNotExist
	ApiRequestLimitError
	ApiRequestError
	ApiParseBodyError
	ApiParseIndexError

	FnError
)
