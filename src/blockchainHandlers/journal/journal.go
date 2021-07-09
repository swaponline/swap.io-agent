package journal

import (
	"swap.io-agent/src/blockchainHandlers"
	"swap.io-agent/src/common/Set"
)

type Journal struct {
	network string
	buf map[string]*blockchainHandlers.SpendsInfo
	spendsAddress  *Set.Set
}

func New(network string) *Journal {
	set := Set.New()
	return &Journal{
		network: network,
		spendsAddress: &set,
	}
}
func (j *Journal) Add(id string, spend blockchainHandlers.Spend) {
	j.spendsAddress.Add(spend.Wallet)

	if _, exist := j.buf[id]; exist {
		j.buf[id].Entries = append(
			j.buf[id].Entries,
			spend,
		)
	} else {
		j.buf[id] = &blockchainHandlers.SpendsInfo{
			Asset: blockchainHandlers.SpendsAsset{
				Id: id,
				Address: id,
				Network: j.network,
				Symbol: j.network + "-" + id,
			},
			Entries: []blockchainHandlers.Spend{spend},
		}
	}
}
func (j *Journal) GetSpends() []blockchainHandlers.SpendsInfo {
	spends := make([]blockchainHandlers.SpendsInfo, len(j.buf))
	for _, spendsInfo := range j.buf {
		spends = append(spends, *spendsInfo)
	}
	return spends
}
func (j *Journal) GetSpendsAddress() []string {
	return j.spendsAddress.Keys()
}