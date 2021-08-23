package geth

import (
	"fmt"
	"net/http"
	"strings"
	"swap.io-agent/src/blockchain/ethereum/api"
)

func (e *Geth) GetBlockTraceByIndex(index string) (interface{}, int) {
	/*
	   curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
	   "method": "debug_traceBlockByNumber",
	   "params": ["0x1f90c0", {"tracer": "callTracer"}],
	   "jsonrpc": "2.0",
	   "id": "2"
	   }'
	*/
	res, err := http.Post(
		e.baseUrl,
		"application/json",
		strings.NewReader(
			fmt.Sprintf(
				`
				{
					"method": "debug_traceBlockByNumber",
					"params": ["%v", {"tracer": "callTracer"}],
					"jsonrpc": "2.0",
					"id": "2"
				}
				`,
				index,
			),
		),
	)
	if err != nil {
		return nil, api.RequestError
	}
	defer res.Body.Close()

	return nil, 0
	//resBodyBytes, err := io.ReadAll(res.Body)
	//if err != nil && err != io.EOF {
	//	log.Println(err)
	//	return nil, api.ParseBodyError
	//}
	//var resBody GetInternalTransactionsResponse
	//if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
	//	log.Println(err, string(resBodyBytes))
	//	return nil, api.RequestError
	//}

	//return &resBody.Result, api.RequestSuccess
}
