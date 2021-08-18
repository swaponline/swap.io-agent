package levelDbStore

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
	"strings"
)

func LinkedListKeyValuesGetCursorData(
	db *leveldb.DB, cursor string,
) ([]string, error) {
	if cursor == "null" {
		return []string{}, nil
	}
	data, err := db.Get([]byte(cursor), nil)
	return strings.Split(string(data), " "), err
}
func LinkedListKeyValuesGetFirstCursor(
	db *leveldb.DB,
	head string,
) (string,string,error) {
	firstCursorDataBytes, err := db.Get([]byte(head), nil)
	if err == leveldb.ErrNotFound {
		return "null","null",nil
	}
	if err != nil {
		return "","",err
	}
	firstCursorData := strings.Split(
		string(firstCursorDataBytes), "|",
	)
	if len(firstCursorData) != 2 {
		return "","",fmt.Errorf("incorrect first cursor data - %v", firstCursorData)
	}
	index, err := strconv.Atoi(firstCursorData[1])
	if err != nil {
		return "","",fmt.Errorf(
			"incorrect first cursor index - %v | %v", index, err,
		)
	}
	if index == 0 {
		return string(firstCursorDataBytes),"null",nil
	}
	return string(firstCursorDataBytes), firstCursorData[0]+"|"+strconv.Itoa(index-1), nil
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
		newTail := []byte(head+`|0`)
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

	keyData := strings.Split(string(tailKey), "|")
	if len(keyData) != 2 {
		return errors.New(
			fmt.Sprintf(
				"invalid tailKey - %v | %v",
				string(tailKey),
				head,
			),
		)
	}

	index, err := strconv.Atoi(keyData[1])
	if err != nil {
		return errors.New(
			fmt.Sprintf(
				"invlid tailKey index - %v",
				keyData[1],
			),
		)
	}

	newTail := []byte(
		head+"|"+strconv.Itoa(index+1),
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
		newTail := []byte(head+"|"+strconv.Itoa(len(data)))
		batch.Put(
			headKey,
			newTail,
		)
		for t:=0; t<len(data); t++ {
			batch.Put(
				[]byte(head+"|"+strconv.Itoa(t)),
				[]byte(data[t]),
			)
		}
		return nil
	}

	keyData := strings.Split(string(tailKey), "|")
	if len(keyData) != 2 {
		return errors.New(
			fmt.Sprintf(
				"invalid tailKey - %v | %v",
				string(tailKey),
				head,
			),
		)
	}

	index, err := strconv.Atoi(keyData[1])
	if err != nil {
		return errors.New(
			fmt.Sprintf(
				"invlid tailKey index - %v",
				keyData[1],
			),
		)
	}

	newTail := []byte(
		head+"|"+strconv.Itoa(index+len(data)),
	)
	batch.Put(
		headKey,
		newTail,
	)
	if err != nil {
		return err
	}
	for t:=0; t<len(data); t++ {
		batch.Put(
			[]byte(head+"|"+strconv.Itoa(index)),
			[]byte(data[t]),
		)
		index+=1
	}

	return err
}
