package handshake

import "swap.io-agent/src/blockchain"

func (a *Api) GetTransactionByHash(
	hash string,
) (
	*blockchain.Transaction,
	error,
) {
	return &blockchain.Transaction{}, nil
}
