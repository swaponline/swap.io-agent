package ethereum

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"strconv"
	"swap.io-agent/src/blockchain"
)

var lastBlockKey = []byte("lastBlock")
var dbpath = "./blockchainIndexes/ethereum"

type BlockchainIndexer struct {
	db              *leveldb.DB
	lastBlock       int
	apiKey          string
	isSynchronize   chan struct{}
	newTransactions chan blockchain.Transaction
}

func InitializeIndexer() *BlockchainIndexer {
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		log.Panicf("db not open %v. err - %v", dbpath, err)
	}

	var lastBlock int
	switch lastBlockStr, err := db.Get(lastBlockKey, nil); err {
		case nil: {
			lastBlockStrNum, err := strconv.Atoi(string(lastBlockStr))
			if err != nil {log.Panicln("num last block not parsed")}
			lastBlock = lastBlockStrNum
		}
		case leveldb.ErrNotFound: {
			err = db.Put(lastBlockKey, []byte("0") , nil)
			if err != nil {log.Panicln("not set value to lastBlockKey")}
			lastBlock = 12833244
		}
		default: log.Panicln("error then get last block index")
	}

	return &BlockchainIndexer{
		db: db,
		lastBlock: lastBlock,
		apiKey: os.Getenv("ETHERSCAN_API_KEY"),
		isSynchronize: make(chan struct{}),
		newTransactions: make(chan blockchain.Transaction),
	}
}