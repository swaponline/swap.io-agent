package blockchain

import "context"

type Broadcaster interface {
	broadcast(hex string) error
}
type Notifier interface {
	notify(ctx context.Context, address string) error
}

type subscribersStore interface {
	GetSubscribersFromAddresses(addresses []string) []string
	SubscribeUserToAddress(userId string, address string) error
	ClearAllUserSubscriptions(userId string) error
}