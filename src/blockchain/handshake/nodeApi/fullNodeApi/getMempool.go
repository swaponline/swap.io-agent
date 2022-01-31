package fullNodeApi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (n *FullNodeApi) GetMempool() ([]string, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(`%v/mempool`, n.baseUrl),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("x", n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("return non ok status code")
	}

	result := make([]string, 0)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
