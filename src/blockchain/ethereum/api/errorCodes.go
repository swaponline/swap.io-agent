package api

const (
	RequestSuccess             = iota
	NotExistBlockError
	RequestLimitError
	RequestError
	ParseBodyError
	ParseIndexError
)