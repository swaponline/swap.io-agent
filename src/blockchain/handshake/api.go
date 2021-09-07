package handshake

import (
	"swap.io-agent/src/blockchain/handshake/nodeApi/fullNodeApi"
)

type Api struct {
	nodeApi *fullNodeApi.FullNodeApi
}

func Test() {
	InitializeApi().nodeApi.GetBlockByIndex(0)
}

func InitializeApi() *Api {
	fullNodeApi := fullNodeApi.InitializeFullNodeApi()

	return &Api{
		nodeApi: fullNodeApi,
	}
}
