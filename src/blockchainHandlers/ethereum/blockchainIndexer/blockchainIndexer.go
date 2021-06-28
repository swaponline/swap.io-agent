package ethereum

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"strconv"
)

var lastBlockKey = []byte("lastBlock")
var dbpath = "./blockchainIndexes/ethereum"

type BlockchainIndexer struct {
	db              *leveldb.DB
	lastBlock       int
	apiKey          string
	isSynchronize   chan struct{}
	newTransactions chan struct{}
}

func InitializeIndexer() *BlockchainIndexer {
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		log.Panicf("db not open %v. err - %v", dbpath, err)
	}

	var lastBlock int
	if lastBlockStr, err := db.Get(lastBlockKey, nil); err == nil {
		lastBlockStrNum, err := strconv.Atoi(string(lastBlockStr))
		if err != nil {log.Panicln("num last block not parsed")}
		lastBlock = lastBlockStrNum
	} else {
		err = db.Put(lastBlockKey, []byte("0") , nil)
		if err != nil {log.Panicln("not set value to lastBlockKey")}
		lastBlock = 0
	}

	return &BlockchainIndexer{
		db: db,
		lastBlock: lastBlock,
		apiKey: os.Getenv("ETHERSCAN_API_KEY"),
		isSynchronize: make(chan struct{}),
		newTransactions: make(chan struct{}),
	}
}