package ethereum

import "swap.io-agent/src/blockchain"

func (a *Api) GetTransactionByHash(
	hash string,
) (
	*blockchain.Transaction,
	error,
) {
	return a.formatter.FormatTransactionFromHash(hash)
}
