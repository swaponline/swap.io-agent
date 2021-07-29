package ethercsan

import "swap.io-agent/src/blockchain/ethereum/api"

type apiError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
type GetTransactionLogsResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  api.TransactionLogs
}
type getCurrentBlockIndexResponse struct {
	Id int `json:"id"`
	Result string `json:"result"`
}
type getBlockResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  api.Block `json:"result"`
}
type getTransactionByHashResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  api.BlockTransaction `json:"result"`
}
