package nodeApi

type Block struct {
	Hash         string        `json:"hash"`
	Height       int           `json:"height"`
	Depth        int           `json:"depth"`
	Version      int           `json:"version"`
	PrevBlock    string        `json:"prevBlock"`
	MerkleRoot   string        `json:"merkleRoot"`
	WitnessRoot  string        `json:"witnessRoot"`
	TreeRoot     string        `json:"treeRoot"`
	ReservedRoot string        `json:"reservedRoot"`
	Time         int           `json:"time"`
	Bits         int           `json:"bits"`
	Nonce        int           `json:"nonce"`
	ExtraNonce   string        `json:"extraNonce"`
	Mask         string        `json:"mask"`
	Txs          []Transaction `json:"txs"`
	Locktime     int           `json:"locktime"`
	Hex          string        `json:"hex"`
}

type Transaction struct {
	Hash        string `json:"hash"`
	Height      int    `json:"height"`
	WitnessHash string `json:"witnessHash"`
	Fee         int    `json:"fee"`
	Rate        int    `json:"rate"`
	Mtime       int    `json:"mtime"`
	Index       int    `json:"index"`
	Version     int    `json:"version"`
	Inputs      []struct {
		Prevout struct {
			Hash  string `json:"hash"`
			Index int    `json:"index"`
		} `json:"prevout"`
		Coin struct {
			Value   int    `json:"value"`
			Address string `json:"address"`
		} `json:"coin"`
	} `json:"inputs"`
	Outputs []struct {
		Value   int    `json:"value"`
		Address string `json:"address"`
	} `json:"outputs"`
	Locktime int    `json:"locktime"`
	Hex      string `json:"hex"`
}
