package levelDbStore

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
	"strings"
)

type writeBuffer struct {
	buf       map[string][]string
	size      int
	lastBlock int
}

type TransactionsStore struct {
	db          *leveldb.DB
	lastBlock   int
	writeBuffer writeBuffer
}
type TransactionsStoreConfig struct {
	Name                 string
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
		writeBuffer: writeBuffer{
			buf: map[string][]string{},
			size: 0,
			lastBlock: -1,
		},
	}, nil
}

func (ts *TransactionsStore) GetLastTransactionBlock() int {
	return ts.lastBlock
}

func (ts *TransactionsStore) GetCursorFromAddress(address string) (string, error) {
	cursor, err := ts.db.Get([]byte(address), nil)
	if err == leveldb.ErrNotFound {
		return "null", nil
	}
	if err != nil {
		return "null", err
	}
	return string(cursor), nil
}
func (ts *TransactionsStore) GetCursorTransactionHashes(
	cursor string,
) (*CursorTransactionHashes, error) {
	hashes := make([]string, 0)
	cursorData, nextCursor, err := LinkedListKeyValuesGetCursorData(
		ts.db, cursor,
	)
	if err != nil {
		return nil, err
	}
	for _, storeData := range cursorData {
		hashTx, _, err := storeDataToTxHashAndBlockIndex(storeData)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hashTx)
	}

	return &CursorTransactionHashes{
		Cursor: cursor,
		NextCursor: nextCursor,
		Hashes: hashes,
	}, nil
}
func (ts *TransactionsStore) GetFirstCursorTransactionHashes(
	address string,
) (*CursorTransactionHashes,error) {
	cursor, _, err := LinkedListKeyValuesGetFirstCursor(
		ts.db, address,
	)
	if err != nil {
		return nil, err
	}
	return ts.GetCursorTransactionHashes(cursor)
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

func txHashAndBlockIndexToStoreData(
	txHash string,
	blockIndex int,
) (string, error) {
	if len(txHash) == 0 {
		return "", fmt.Errorf(
			"incorrect hashTx - %v | block index - %v",
			txHash,
			blockIndex,
		)
	}
	return txHash+"|"+strconv.Itoa(blockIndex), nil
}
func storeDataToTxHashAndBlockIndex(storeData string) (string, int, error) {
	hashTxAndBlockIndex := strings.Split(storeData, "|")
	if len(hashTxAndBlockIndex) != 2 || len(hashTxAndBlockIndex[0]) == 0 {
		return "", 0, fmt.Errorf("inccorrect storeData - %v", storeData)
	}
	blockIndex, err := strconv.Atoi(hashTxAndBlockIndex[1])
	if err != nil {
		return "", 0, fmt.Errorf("inccorrect storeData - %v", storeData)
	}
	return hashTxAndBlockIndex[0], blockIndex, nil
}

func (ts *TransactionsStore) WriteLastIndexedTransactions(
	AddressHashTransactions map[string][]string,
	indexBlock int,
) error {
	for address, hashes := range AddressHashTransactions {
		hashIndexTransactions := make([]string, 0)
		for _, hash := range hashes {
			storeData, err := txHashAndBlockIndexToStoreData(hash, indexBlock)
			if err != nil {
				return err
			}
			hashIndexTransactions = append(
				hashIndexTransactions,
				storeData,
			)
		}

		ts.writeBuffer.buf[address] = append(
			ts.writeBuffer.buf[address],
			hashIndexTransactions...
		)
		ts.writeBuffer.size+=len(hashes)
	}

	ts.lastBlock             = indexBlock
	ts.writeBuffer.lastBlock = indexBlock
	if ts.writeBuffer.size > 1024 {
		return ts.Flush()
	}

	return nil
}
func (ts *TransactionsStore) Flush() error {
	if ts.writeBuffer.size > 0 ||
	   ts.writeBuffer.lastBlock != -1 {
		dbTransaction, err := ts.db.OpenTransaction()
		if err != nil {
			return err
		}

		batch := new(leveldb.Batch)
		for address, hashIndexTransactions := range ts.writeBuffer.buf {
			err := LinkedListKeyValuesPush(
				ts.db,
				batch,
				address,
				hashIndexTransactions,
			)
			if err != nil {
				dbTransaction.Discard()
				return err
			}
		}

		err = dbTransaction.Write(batch, nil)
		if err != nil {
			dbTransaction.Discard()
			return err
		}

		err = dbTransaction.Put(
			lastBlockKey,
			[]byte(strconv.Itoa(ts.writeBuffer.lastBlock)),
			nil,
		)
		if err != nil {
			dbTransaction.Discard()
			return err
		}

		err = dbTransaction.Commit()
		if err != nil {
			dbTransaction.Discard()
			return err
		}
		log.Println("written transactions -", ts.writeBuffer.size)
		log.Println("updated lastBlock -", ts.writeBuffer.lastBlock)

		ts.writeBuffer.buf       = make(map[string][]string)
		ts.writeBuffer.size      = 0
		ts.writeBuffer.lastBlock = -1
	}

	return nil
}

func (ts *TransactionsStore) Start() {}
func (ts *TransactionsStore) Stop() error {
	return nil
}
func (ts *TransactionsStore) Status() error {
	return nil
}