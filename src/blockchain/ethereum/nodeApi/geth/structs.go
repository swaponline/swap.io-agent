package geth

import "swap.io-agent/src/blockchain/ethereum/nodeApi"

type apiError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
type GetTransactionLogsResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  nodeApi.TransactionLogs
}
type GetInternalTransactionsResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  nodeApi.InteranlTransaction
}
type GetCurrentBlockIndexResponse struct {
	Id     int    `json:"id"`
	Result string `json:"result"`
}
type getBlockResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Result  nodeApi.Block `json:"result"`
}
type getTransactionByHashResponse struct {
	Jsonrpc string                   `json:"jsonrpc"`
	Id      int                      `json:"id"`
	Result  nodeApi.BlockTransaction `json:"result"`
}
