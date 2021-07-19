package levelDbStore

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
)

type TransactionStore struct {
	db *leveldb.DB
	lastBlock int
}
type TransactionStoreConfig struct {
	name string
	defaultScannedBlocks int
}

const dbDir = "./blockchainIndexes/"
var lastBlockKey = []byte("lastBlock")

func InitialiseTransactionStore(config TransactionStoreConfig) *TransactionStore {
	db, err := leveldb.OpenFile(dbDir+config.name, nil)
	if err != nil {
		log.Panicf("db not open %v. err - %v", dbDir+config.name, err)
	}

	var lastBlock int
	switch lastBlockStr, err := db.Get(lastBlockKey, nil); err {
	case nil: {
		lastBlockStrNum, err := strconv.Atoi(string(lastBlockStr))
		if err != nil {log.Panicln("num last block not parsed")}
		lastBlock = lastBlockStrNum
	}
	case leveldb.ErrNotFound: {
		err = db.Put(
			lastBlockKey,
			[]byte(strconv.Itoa(config.defaultScannedBlocks)),
			nil,
		)
		if err != nil {log.Panicln("not set value to lastBlockKey")}
		lastBlock = config.defaultScannedBlocks
	}
	default: log.Panicln("error then get last block index")
	}

	return &TransactionStore{
		db: db,
		lastBlock: lastBlock,
	}
}
func (ts *TransactionStore) WriteBlockIndexedTransactions() error {
	return nil
}

func (ts *TransactionStore) Start() {}
func (ts *TransactionStore) Stop() error {
	return nil
}
func (ts *TransactionStore) Status() error {
	return nil
}