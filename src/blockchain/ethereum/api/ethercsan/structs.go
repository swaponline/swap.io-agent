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
	Result  struct {
		BlockHash         string      `json:"blockHash"`
		BlockNumber       string      `json:"blockNumber"`
		ContractAddress   string      `json:"contractAddress"`
		CumulativeGasUsed string      `json:"cumulativeGasUsed"`
		From              string      `json:"from"`
		GasUsed           string      `json:"gasUsed"`
		Logs              []api.TransactionLog `json:"logs"`
		LogsBloom        string `json:"logsBloom"`
		Status           string `json:"status"`
		To               string `json:"to"`
		TransactionHash  string `json:"transactionHash"`
		TransactionIndex string `json:"transactionIndex"`
		Type             string `json:"type"`
	} `json:"result"`
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
