package fullNodeApi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (n *FullNodeApi) PushTx(hex string) ([]byte, error) {
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

	result, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("return non ok status code")
	}

	return result, nil
}
