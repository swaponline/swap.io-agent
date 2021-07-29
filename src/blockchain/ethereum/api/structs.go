package api

type Block struct {
	Difficulty       string `json:"difficulty"`
	ExtraData        string `json:"extraData"`
	GasLimit         string `json:"gasLimit"`
	GasUsed          string `json:"gasUsed"`
	Hash             string `json:"hash"`
	LogsBloom        string `json:"logsBloom"`
	Miner            string `json:"miner"`
	MixHash          string `json:"mixHash"`
	Nonce            string `json:"nonce"`
	Number           string `json:"number"`
	ParentHash       string `json:"parentHash"`
	ReceiptsRoot     string `json:"receiptsRoot"`
	Sha3Uncles       string `json:"sha3Uncles"`
	Size             string `json:"size"`
	StateRoot        string `json:"stateRoot"`
	Timestamp        string `json:"timestamp"`
	TotalDifficulty  string `json:"totalDifficulty"`
	Transactions     []BlockTransaction
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []interface{} `json:"uncles"`
}
type BlockTransaction struct {
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	From              string `json:"from"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	Hash              string `json:"hash"`
	Input             string `json:"input"`
	Nonce             string `json:"nonce"`
	To                string `json:"to"`
	TransactionIndex  string `json:"transactionIndex"`
	Value             string `json:"value"`
	Type              string `json:"type"`
	V                 string `json:"v"`
	R                 string `json:"r"`
	S                 string `json:"s"`
	AllSpendAddresses []string
}
type TransactionLogs struct {
	BlockHash         string      `json:"blockHash"`
	BlockNumber       string      `json:"blockNumber"`
	ContractAddress   string      `json:"contractAddress"`
	CumulativeGasUsed string      `json:"cumulativeGasUsed"`
	From              string      `json:"from"`
	GasUsed           string      `json:"gasUsed"`
	Logs              []TransactionLog `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Status            string `json:"status"`
	To                string `json:"to"`
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
	Type              string `json:"type"`
}
type TransactionLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}
