package redisStore

type ISubscribersStore interface {
	GetSubscribersFromAddresses(addresses []string) []string
	AddSubscription(userId string, address string) error
	RemoveSubscription(userId string, address string) error
	RemoveSubscriptions(userId string) error
}
