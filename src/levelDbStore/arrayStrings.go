package levelDbStore

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"strings"
)

type GetPutLevelDb interface {
	Get([]byte, *opt.ReadOptions) ([]byte, error)
	Put([]byte, []byte, *opt.WriteOptions) error
}

func ArrayStringPush(db GetPutLevelDb, key string, value []string) error {
	keyValue, err := db.Get([]byte(key), nil)
	if err != leveldb.ErrNotFound {
		return err
	}

	newKeyValue := strings.Join(
		append(
			strings.Split(string(keyValue), " "),
			value...
		),
		" ",
	)

	return db.Put(
		[]byte(key),
		[]byte(newKeyValue),
		nil,
	)
}