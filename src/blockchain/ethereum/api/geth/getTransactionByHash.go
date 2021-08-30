package geth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"swap.io-agent/src/blockchain/ethereum/api"
)

/*
curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
"method": "eth_getTransactionByHash",
"params": ["0x3837729b0628eef5acfab9a726271b1f2a9365e59a29f889e73c663d36088be3"],
"jsonrpc": "2.0",
"id": "2"
}'
*/
func (e *Geth) GetTransactionByHash(
	hash string,
) (*api.BlockTransaction, int) {

	res, err := http.Post(
		e.baseUrl,
		"application/json",
		strings.NewReader(
			fmt.Sprintf(
				`{
					"jsonrpc":"2.0",
					"method":"eth_getTransactionByHash",
					"params":["%v"],
					"id":1
				}`,
				hash,
			),
		),
	)
	if err != nil {
		return nil, api.RequestError
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		return nil, api.ParseBodyError
	}

	var resBody getTransactionByHashResponse
	if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, api.RequestError
	}

	return &resBody.Result, api.RequestSuccess
}
