package handshake

func (a *Api) PushTx(hex string) ([]byte, error) {
	return a.nodeApi.PushTx(hex)
}
