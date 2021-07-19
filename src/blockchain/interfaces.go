package blockchain

import "context"

type Broadcaster interface {
	broadcast(hex string) error
}
type Notifier interface {
	notify(ctx context.Context, address string) error
}

