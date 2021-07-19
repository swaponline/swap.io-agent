package redisStore

type SubscribersStore interface {
	GetSubscribersFromAddresses(addresses []string) []string
	SubscribeUserToAddress(userId string, address string) error
	ClearAllUserSubscriptions(userId string) error
}