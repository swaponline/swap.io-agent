package journal

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/common/Set"
)

type Journal struct {
	network       string
	buf           map[string]*blockchain.SpendsInfo
	spendsAddress *Set.Set
}

func New(network string) *Journal {
	set := Set.New()
	return &Journal{
		network:       network,
		spendsAddress: &set,
		buf:           make(map[string]*blockchain.SpendsInfo),
	}
}
func (j *Journal) Add(id string, spend blockchain.Spend) {
	if len(spend.Wallet) != 0 {
		j.spendsAddress.Add(spend.Wallet)
	}

	if _, exist := j.buf[id]; exist {
		j.buf[id].Entries = append(
			j.buf[id].Entries,
			spend,
		)
	} else {
		j.buf[id] = &blockchain.SpendsInfo{
			Asset: blockchain.SpendsAsset{
				Id:      id,
				Address: id,
				Network: j.network,
				Symbol:  j.network + "-" + id,
			},
			Entries: []blockchain.Spend{spend},
		}
	}
}
func (j *Journal) GetSpends() []blockchain.SpendsInfo {
	spends := make([]blockchain.SpendsInfo, 0)
	for _, spendsInfo := range j.buf {
		spends = append(spends, *spendsInfo)
	}
	return spends
}
func (j *Journal) GetSpendsAddress() []string {
	return j.spendsAddress.Keys()
}
