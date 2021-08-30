package levelDbStore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

func dataToKey(head string, index int) string {
	return head + "|" + strconv.Itoa(index)
}
func keyToData(key string) (string, int, error) {
	parsedKey := strings.Split(key, "|")
	if len(parsedKey) != 2 {
		return "null", -1, fmt.Errorf(
			"key - %v",
			key,
		)
	}
	index, err := strconv.Atoi(parsedKey[1])
	if err != nil {
		return "null", -1, fmt.Errorf(
			"incorrect key index - %v", parsedKey[1],
		)
	}
	return parsedKey[0], index, nil
}

func LinkedListKeyValuesGetCursorData(
	db *leveldb.DB,
	cursor string,
) ([]string, string, error) {
	if cursor == "null" {
		return []string{}, "null", nil
	}
	_, _, err := keyToData(cursor)
	if err != nil {
		return []string{}, "null", nil
	}

	data, err := db.Get([]byte(cursor), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return []string{}, "null", err
	}

	nextCursor, err := LinkedListKeyValuesGetNextCursor(cursor)
	if err != nil {
		return []string{}, "null", nil
	}

	return strings.Split(string(data), " "), nextCursor, nil
}
func LinkedListKeyValuesGetFirstCursor(
	db *leveldb.DB, head string,
) (string, string, error) {
	if head == "null" {
		return "null", "null", nil
	}

	cursor, err := db.Get([]byte(head), nil)
	if err == leveldb.ErrNotFound {
		return "null", "null", nil
	}
	if err != nil {
		return "null", "null", err
	}

	nextCursor, err := LinkedListKeyValuesGetNextCursor(string(cursor))
	if err != nil {
		return "null", "null", err
	}

	return string(cursor), nextCursor, nil
}

func LinkedListKeyValuesGetNextCursor(
	cursor string,
) (string, error) {
	if cursor == "null" {
		return "null", nil
	}

	head, index, err := keyToData(cursor)
	if err != nil {
		return "null", err
	}
	if index <= 0 {
		return "null", nil
	}
	return dataToKey(head, index-1), nil
}

func LinkedListKeyValuesPush(
	db *leveldb.DB,
	batch *leveldb.Batch,
	head string,
	data []string,
) error {
	if len(data) == 0 {
		return nil
	}
	headKey := []byte(head)
	tailKey, err := db.Get(headKey, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	if err == leveldb.ErrNotFound {
		newTail := []byte(head + `|0`)
		batch.Put(
			headKey,
			newTail,
		)
		batch.Put(
			newTail,
			[]byte(strings.Join(data, " ")),
		)
		return nil
	}

	_, index, err := keyToData(string(tailKey))
	if err != nil {
		return err
	}

	newTail := []byte(
		head + "|" + strconv.Itoa(index+1),
	)
	batch.Put(
		headKey,
		newTail,
	)
	if err != nil {
		return err
	}
	batch.Put(
		newTail,
		[]byte(strings.Join(data, " ")),
	)

	return err
}
func LinkedListKeyValuePush(
	db *leveldb.DB,
	batch *leveldb.Batch,
	head string,
	data []string,
) error {
	if len(data) == 0 {
		return nil
	}
	headKey := []byte(head)
	tailKey, err := db.Get(headKey, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	if err == leveldb.ErrNotFound {
		newTail := []byte(head + "|" + strconv.Itoa(len(data)))
		batch.Put(
			headKey,
			newTail,
		)
		for t := 0; t < len(data); t++ {
			batch.Put(
				[]byte(head+"|"+strconv.Itoa(t)),
				[]byte(data[t]),
			)
		}
		return nil
	}

	_, index, err := keyToData(string(tailKey))
	if err != nil {
		return err
	}

	newTail := []byte(
		head + "|" + strconv.Itoa(index+len(data)),
	)
	batch.Put(
		headKey,
		newTail,
	)
	if err != nil {
		return err
	}
	for t := 0; t < len(data); t++ {
		batch.Put(
			[]byte(head+"|"+strconv.Itoa(index)),
			[]byte(data[t]),
		)
		index += 1
	}

	return err
}
