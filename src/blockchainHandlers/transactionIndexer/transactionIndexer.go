package transactionIndexer

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

type TransactionIndexer struct {}
type Options struct {
	DbPath string
	FirstBlock string
}

func New(options Options) {
	db, err := leveldb.OpenFile(options.DbPath, nil)
	if err != nil {
		log.Panicf("db not open %v. err - %v", options.DbPath, err)
	}
	lastBlock, err := db.Get(lastBlockKey, nil)
	//last block not set
	if err != nil {
		lastBlock = []byte(options.FirstBlock)
		db.Put(lastBlockKey, lastBlock, nil)
	}

	log.Println(lastBlock)
}

var lastBlockKey = []byte("lastBlock")