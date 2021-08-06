package levelDbStore

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"strings"
)

type IGetPutLevelDb interface {
	Has([]byte, *opt.ReadOptions) (bool, error)
	Get([]byte, *opt.ReadOptions) ([]byte, error)
	Put([]byte, []byte, *opt.WriteOptions) error
}

func ArrayStringPush(db IGetPutLevelDb, key string, values []string) error {
	if len(values) == 0 {
		return nil
	}

	keyByte := []byte(key)

	keyValue, err := db.Get(keyByte, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}

	var newKeyValue []byte
	if len(keyValue) == 0 {
		newKeyValue = []byte(strings.Join(values, " "))
	} else {
		newKeyValue = keyValue
		for t:=0; t<len(values); t++ {
			newKeyValue = append(
				append(newKeyValue, ' '),
				[]byte(values[t])...
			)
		}
	}

	return db.Put(
		[]byte(key),
		newKeyValue,
		nil,
	)
}