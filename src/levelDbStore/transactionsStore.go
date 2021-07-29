package levelDbStore

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
	"strings"
)

type TransactionsStore struct {
	db *leveldb.DB
	lastBlock int
}
type TransactionsStoreConfig struct {
	Name string
	DefaultScannedBlocks int
}

const dbDir = "./blockchainIndexes/"
var lastBlockKey = []byte("lastBlock")

func InitialiseTransactionStore(config TransactionsStoreConfig) (*TransactionsStore, error) {
	db, err := leveldb.OpenFile(dbDir+config.Name, nil)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("db not open %v. err - %v", dbDir+config.Name, err),
		)
	}

	var lastBlock int
	switch lastBlockStr, err := db.Get(lastBlockKey, nil); err {
	case nil: {
		lastBlockStrNum, err := strconv.Atoi(string(lastBlockStr))
		if err != nil {
			return nil, errors.New("num last block not parsed")
		}
		lastBlock = lastBlockStrNum
	}
	case leveldb.ErrNotFound: {
		err = db.Put(
			lastBlockKey,
			[]byte(strconv.Itoa(config.DefaultScannedBlocks)),
			nil,
		)
		if err != nil {
			return nil, errors.New("not set value to lastBlockKey")
		}
		lastBlock = config.DefaultScannedBlocks
	}
	default: return nil, errors.New("error then get last block index")
	}

	return &TransactionsStore{
		db: db,
		lastBlock: lastBlock,
	}, nil
}

func (ts *TransactionsStore) GetLastTransactionBlock() int {
	return ts.lastBlock
}
func (ts *TransactionsStore) GetAddressTransactionsHash(
	address string,
	startTime int,
	endTime int,
) ([]string, error) {
	transactionsInfoBytes, err := ts.db.Get([]byte(address), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return make([]string, 0), nil
		}
		return nil, err
	}

	transactionsInfo := strings.Split(
		string(transactionsInfoBytes),
		" ",
	)
	transactionsHash := make([]string, 0)
	for _, transactionInfo := range transactionsInfo {
		hashTimestamp := strings.Split(transactionInfo, "|")
		if len(hashTimestamp) != 2 {
			return nil, errors.New("not get hash|time transaction info from db")
		}
		timestamp, err := strconv.Atoi(hashTimestamp[1])
		if err != nil {
			return nil, err
		}

		if startTime <= timestamp && timestamp <= endTime {
			transactionsHash = append(transactionsHash, hashTimestamp[0])
		}
	}

	return transactionsHash, nil
}
func (ts *TransactionsStore) WriteLastIndexedBlockTransactions(
	indexedTransactions *map[string][]string,
	indexBlock int,
	timestampBlock int,
) error {
	bdTransaction, err := ts.db.OpenTransaction()
	if err != nil {
		return err
	}
	log.Println(indexedTransactions)
	for address, transactions := range *indexedTransactions {
		//format transactions
		formattedTransactions := make([]string, len(transactions))
		for index, hashTransaction := range transactions {
			formattedTransactions[index] = fmt.Sprintf(
				`%v|%v`, hashTransaction, timestampBlock,
			)
		}
		// push back address transaction|timestampBlock
		err = ArrayStringPush(
			bdTransaction, address, formattedTransactions,
		)
		if err != nil {
			bdTransaction.Discard()
			return err
		}
	}
	// update lastBlock
	err = bdTransaction.Put(
		lastBlockKey,
		[]byte(strconv.Itoa(indexBlock)),
		nil,
	)
	if err != nil {
		bdTransaction.Discard()
		return err
	}

	// commit transaction
	err = bdTransaction.Commit()
	if err != nil {
		bdTransaction.Discard()
		return err
	}

	ts.lastBlock = indexBlock

	return nil
}

func (ts *TransactionsStore) Start() {}
func (ts *TransactionsStore) Stop() error {
	return nil
}
func (ts *TransactionsStore) Status() error {
	return nil
}