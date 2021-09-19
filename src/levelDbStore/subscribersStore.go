package levelDbStore

import (
	"fmt"
	"log"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type SubscribersStore struct {
	db *leveldb.DB
}

const subscribersStoreDbDir = "./subscriptions"

func InitialiseSubscriberStore() (*SubscribersStore, error) {
	db, err := leveldb.OpenFile(subscribersStoreDbDir, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"db not open %v. err - %v", subscribersStoreDbDir, err,
		)
	}

	return &SubscribersStore{
		db: db,
	}, nil
}

func (s *SubscribersStore) GetSubscriptions(userId string) ([]string, error) {
	subscriptions := []string{}
	iter := s.db.NewIterator(util.BytesPrefix([]byte(userId+"|")), nil)
	for iter.Next() {
		subscriptionData := strings.Split(string(iter.Key()), "|")
		if len(subscriptionData) != 2 {
			log.Printf("incorrect subscriptionData levelDb %#v", subscriptionData)
			return subscriptions, fmt.Errorf(
				"incorrect subscriptionData levelDb %#v",
				subscriptionData,
			)
		}
		subscriptions = append(subscriptions, subscriptionData[1])
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return subscriptions, err
	}

	return subscriptions, nil
}

func (s *SubscribersStore) AddSubscription(
	userId string, address string,
) error {
	return s.db.Put([]byte(userId+"|"+address), []byte{}, nil)
}
func (s *SubscribersStore) RemoveSubscription(
	userId string, address string,
) error {
	return s.db.Delete([]byte(userId+"|"+address), nil)
}

func (*SubscribersStore) Start() {}
func (*SubscribersStore) Stop() error {
	return nil
}
func (*SubscribersStore) Status() error {
	return nil
}
