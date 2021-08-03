package levelDbStore

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"strings"
)

type IGetPutLevelDb interface {
	Get([]byte, *opt.ReadOptions) ([]byte, error)
	Put([]byte, []byte, *opt.WriteOptions) error
}

func ArrayStringPush(db IGetPutLevelDb, key string, value []string) error {
	keyValue, err := db.Get([]byte(key), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}

	var newKeyValue string
	if len(keyValue) == 0 {
		newKeyValue = strings.Join(value, " ")
	} else {
		newKeyValue = strings.Join(
			append(
				strings.Split(string(keyValue), " "),
				value...
			),
			" ",
		)
	}


	return db.Put(
		[]byte(key),
		[]byte(newKeyValue),
		nil,
	)
}