package blockchainHandlers

type Journal struct {
	network string
	buf map[string]*SpendsInfo
}

func New(network string) *Journal {
	return &Journal{
		network: network,
	}
}
func (j *Journal) Add(id string, spend Spend)  {
	if _, exist := j.buf[id]; exist {
		j.buf[id].Entries = append(
			j.buf[id].Entries,
			spend,
		)
	} else {
		j.buf[id] = &SpendsInfo{
			Asset: SpendsAsset{
				Id: id,
				Address: id,
				Network: j.network,
				Symbol: j.network + "-" + id,
			},
			Entries: []Spend{spend},
		}
	}
}
func (j *Journal) GetSpends() []SpendsInfo {
	spends := make([]SpendsInfo, len(j.buf))
	for _, spendsInfo := range j.buf {
		spends = append(spends, *spendsInfo)
	}
	return spends
}