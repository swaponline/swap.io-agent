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

func LinkedListKeyValuesGetNextCursor(
	head string,
	current int,
	size int,
) string {
	if current >= size-1 {
		return "null"
	}

	return dataToKey(head, current+1)
}

func LinkedListKeyValuesGetFirstCursor(
	db *leveldb.DB, head string,
) (string, string, error) {
	if head == "null" {
		return "null", "null", nil
	}

	sizeBytes, err := db.Get([]byte(head), nil)
	if err == leveldb.ErrNotFound {
		return "null", "null", nil
	}
	if err != nil {
		return "null", "null", err
	}
	size, err := strconv.Atoi(string(sizeBytes))
	if err != nil {
		return "null", "null", err
	}

	nextCursor := LinkedListKeyValuesGetNextCursor(head, 0, size)
	if err != nil {
		return "null", "null", err
	}

	return dataToKey(head, 0), nextCursor, nil
}
func LinkedListKeyValuesGetCursorData(
	db *leveldb.DB,
	cursor string,
) ([]string, string, error) {
	if cursor == "null" {
		return []string{}, "null", nil
	}
	head, current, err := keyToData(cursor)
	if err != nil {
		return []string{}, "null", nil
	}
	sizeBytes, err := db.Get([]byte(head), nil)
	if err == leveldb.ErrNotFound {
		return []string{}, "null", err
	}
	if err != nil {
		return []string{}, "null", err
	}
	size, err := strconv.Atoi(string(sizeBytes))
	if err != nil {
		return []string{}, "null", err
	}

	data, err := db.Get([]byte(cursor), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return []string{}, "null", err
	}

	nextCursor := LinkedListKeyValuesGetNextCursor(head, current, size)
	if err != nil {
		return []string{}, "null", nil
	}

	return strings.Split(string(data), " "), nextCursor, nil
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
	sizeBytes, err := db.Get(headKey, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	if err == leveldb.ErrNotFound {
		newTail := []byte(head + `|0`)
		batch.Put(
			headKey,
			[]byte(`1`),
		)
		batch.Put(
			newTail,
			[]byte(strings.Join(data, " ")),
		)
		return nil
	}

	size, err := strconv.Atoi(string(sizeBytes))
	if err != nil {
		return err
	}
	newSize := size + 1
	newTail := []byte(
		head + "|" + strconv.Itoa(size),
	)
	batch.Put(
		headKey,
		[]byte(strconv.Itoa(newSize)),
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
