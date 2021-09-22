package socketServer

type SubscriptionEventPayload struct {
	Address string `json:"address"`
}
type SubscriptionsSize struct {
	Size int `json:"size"`
}
