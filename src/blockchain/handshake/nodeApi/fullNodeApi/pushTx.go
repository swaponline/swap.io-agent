package fullNodeApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (n *FullNodeApi) PushTx(hex string) (interface{}, error) {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(`%v/broadcast`, n.baseUrl),
		strings.NewReader(fmt.Sprintf(`{"tx":"%v"}`, hex)),
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("x", n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
