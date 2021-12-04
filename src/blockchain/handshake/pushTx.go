package handshake

func (a *Api) PushTx(hex string) (interface{}, error) {
	return a.nodeApi.PushTx(hex)
}
