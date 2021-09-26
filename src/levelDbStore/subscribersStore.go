package levelDbStore

import (
	"fmt"
	"log"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"swap.io-agent/src/env"
)

type SubscribersStore struct {
	db *leveldb.DB
}

const subscribersStoreDbDir = "./subscriptions"

func InitialiseSubscribersStore() (*SubscribersStore, error) {
	db, err := leveldb.OpenFile(subscribersStoreDbDir+"/"+env.BLOCKCHAIN, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"db not open %v. err - %v", subscribersStoreDbDir, err,
		)
	}

	return &SubscribersStore{
		db: db,
	}, nil
}

func createSubsctiptin(userId, address string) string {
	return userId + "|" + address
}
func parseSubscription(subscription string) (string, string) {
	subscriptionData := strings.Split(subscription, "|")
	if len(subscriptionData) != 2 {
		log.Panicf("incorrect subscriptionData levelDb %#v", subscriptionData)
	}
	return subscriptionData[0], subscriptionData[1]
}

func (s *SubscribersStore) GetAllSubscriptions() (map[string][]string, error) {
	subscriptions := make(map[string][]string)

	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		userId, address := parseSubscription(string(iter.Key()))
		subscriptions[userId] = append(subscriptions[userId], address)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return subscriptions, err
	}

	return subscriptions, nil
}
func (s *SubscribersStore) GetSubscriptions(userId string) ([]string, error) {
	subscriptions := []string{}
	iter := s.db.NewIterator(util.BytesPrefix([]byte(userId+"|")), nil)
	for iter.Next() {
		_, address := parseSubscription(string(iter.Key()))
		subscriptions = append(subscriptions, address)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return subscriptions, err
	}

	return subscriptions, nil
}
func (s *SubscribersStore) GetSubscriptionsSize(userId string) (int, error) {
	subscriptionsSize := 0
	iter := s.db.NewIterator(util.BytesPrefix([]byte(userId+"|")), nil)
	for iter.Next() {
		subscriptionsSize += 1
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return 0, err
	}

	return subscriptionsSize, nil
}

func (s *SubscribersStore) AddSubscription(
	userId string, address string,
) error {
	return s.db.Put(
		[]byte(createSubsctiptin(userId, address)),
		[]byte{},
		nil,
	)
}
func (s *SubscribersStore) RemoveSubscription(
	userId string, address string,
) error {
	return s.db.Delete(
		[]byte(createSubsctiptin(userId, address)),
		nil,
	)
}

func (*SubscribersStore) Start() {}
func (*SubscribersStore) Stop() error {
	return nil
}
func (*SubscribersStore) Status() error {
	return nil
}
