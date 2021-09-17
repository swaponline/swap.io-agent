package levelDbStore

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type SubscribersStore struct {
	db *leveldb.DB
}

const subscribersStoreDbDir = "./subscribers"

func InitialiseSubscriberStore(
	config TransactionsStoreConfig,
) (*SubscribersStore, error) {
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
func (s *SubscribersStore) GetCountSubcribers(userId string) (int, error) {
	keys := 0
	iter := s.db.NewIterator(util.BytesPrefix([]byte(userId+"|")), nil)
	for iter.Next() {
		keys += 1
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return 0, err
	}

	return keys, nil
}

func (*SubscribersStore) Start() {}
func (*SubscribersStore) Stop() error {
	return nil
}
func (*SubscribersStore) Status() error {
	return nil
}
