package redisStore

type ISubscribersStore interface {
	GetSubscribersFromAddresses(addresses []string) []string
	SubscribeUserToAddress(userId string, address string) error
	ClearAllUserSubscriptions(userId string) error
}