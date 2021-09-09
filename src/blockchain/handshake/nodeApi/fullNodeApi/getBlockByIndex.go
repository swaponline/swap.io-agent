package fullNodeApi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"swap.io-agent/src/blockchain"
	nodeApi "swap.io-agent/src/blockchain/handshake/nodeApi"
)

func (n FullNodeApi) GetBlockByIndex(index int) (*nodeApi.Block, int) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%v/block/%v", n.baseUrl, index),
		nil,
	)
	if err != nil {
		log.Println(err)
		return nil, blockchain.ApiRequestError
	}
	req.SetBasicAuth("x", n.apiKey)

	resp, err := n.client.Do(req)
	if resp.StatusCode == http.StatusNotFound {
		return nil, blockchain.ApiNotExistBlockError
	}
	if err != nil {
		log.Println(err)
		return nil, blockchain.ApiRequestError
	}
	defer resp.Body.Close()

	var block *nodeApi.Block
	if err := json.NewDecoder(resp.Body).Decode(&block); err != nil {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}

	return block, blockchain.ApiRequestSuccess
}
