package geth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
)

/*
curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
"method": "debug_traceTransaction",
"params": ["0x3837729b0628eef5acfab9a726271b1f2a9365e59a29f889e73c663d36088be3", {"tracer": "callTracer"}],
"jsonrpc": "2.0",
"id": "2"
}'
*/
func (e *Geth) GetInternalTransaction(hash string) (*nodeApi.InteranlTransaction, int) {
	res, err := http.Post(
		e.baseUrl,
		"application/json",
		strings.NewReader(
			fmt.Sprintf(
				`
				{
					"method": "debug_traceTransaction",
					"params": ["%v", {"tracer": "callTracer"}],
					"jsonrpc": "2.0",
					"id": "2"
				}
				`,
				hash,
			),
		),
	)
	if err != nil {
		return nil, blockchain.ApiRequestError
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}
	var resBody GetInternalTransactionsResponse
	if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, blockchain.ApiParseBodyError
	}

	return &resBody.Result, blockchain.ApiRequestSuccess
}
