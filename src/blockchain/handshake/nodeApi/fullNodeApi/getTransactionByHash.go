package fullNodeApi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
)

func (n FullNodeApi) GetTransactionByHash(hash string) (*nodeApi.Transaction, int) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%v/tx/%v", n.baseUrl, hash),
		nil,
	)
	if err != nil {
		log.Println(err)
		return nil, blockchain.ApiRequestError
	}
	req.SetBasicAuth("x", n.apiKey)

	resp, err := n.client.Do(req)
	if resp.StatusCode == http.StatusNotFound {
		return nil, blockchain.ApiNotExist
	}
	if err != nil {
		log.Println(err)
		return nil, blockchain.ApiRequestError
	}
	defer resp.Body.Close()

	var transaction *nodeApi.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}
	log.Println(transaction)

	return transaction, blockchain.ApiRequestSuccess
}
