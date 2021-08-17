package blockchain

type Transaction struct {
	From              string       `json:"from"`
	Gas               string       `json:"gas"`
	Hash              string       `json:"hash"`
	Nonce             string       `json:"nonce"`
	To                string       `json:"to"`
	TransactionIndex  string       `json:"transactionIndex"`
	Value             string       `json:"value"`
	GasUsed           string       `json:"gas_used"`
	BlockHash         string       `json:"block_hash"`
	BlockNumber       string       `json:"block_number"`
	GasPrice          string       `json:"gas_price"`
	Timestamp         string       `json:"timestamp"`
	BlockMiner        string       `json:"block_miner"`
	Journal           []SpendsInfo `json:"journal"`
	AllSpendAddresses []string
}
type SpendsInfo struct {
	Asset     SpendsAsset `json:"asset"`
	Entries   []Spend     `json:"entries"`
}
type SpendsAsset struct {
	Id      string `json:"id"`
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
	Network string `json:"network"`
}
type Spend struct {
	Wallet string `json:"wallet"`
	Value  string `json:"value"`
	Label  string `json:"label,omitempty"`
}
type TransactionPipeData struct {
	Subscribers []string
	Transaction *Transaction
}