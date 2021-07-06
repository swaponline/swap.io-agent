package blockchainHandlers

type Transaction struct {
	From             string `json:"from"`
	Gas              int    `json:"gas"`
	Hash             string `json:"hash"`
	Nonce            int    `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex int    `json:"transactionIndex"`
	Value            int  `json:"value"`
	GasUsed          int    `json:"gas_used"`
	BlockHash        string `json:"block_hash"`
	BlockNumber      int    `json:"block_number"`
	GasPrice         int  `json:"gas_price"`
	Timestamp        int    `json:"timestamp"`
	BlockMiner       string `json:"block_miner"`
	Journal          [] struct {
		Asset struct{
			Symbol   string `json:"symbol"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
		} `json:"asset"`
		Entries []struct {
			Wallet string `json:"wallet"`
			Value  int64  `json:"value"`
		} `json:"entries"`
	} `json:"journal"`

	AllSpendAddresses []string
}
