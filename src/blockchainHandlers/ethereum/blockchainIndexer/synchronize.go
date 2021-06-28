package ethereum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type etherscanCurrentBlockRes struct {
	Id int `json:"id"`
	Result string `json:"result"`
}

func (indexer *BlockchainIndexer) Synchronize() {
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=%v",
			indexer.apiKey,
		),
	)
	if err != nil {log.Panicf("not get last block err %v", err)}

	var currentBlockInfo etherscanCurrentBlockRes
	if err = json.NewDecoder(res.Body).Decode(&currentBlockInfo); err != nil {
		log.Panicf("not get last block err %v", err)
	}

	currentBlockId, err := strconv.ParseInt(
		strings.TrimPrefix(currentBlockInfo.Result, "0x"),
		16,
		64,
	)
	if err != nil {log.Panicf("currentBlockId not parse to int %v", err)}

	log.Println(currentBlockId, "current ethereum block")

	close(indexer.isSynchronize)
}