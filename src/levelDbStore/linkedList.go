package levelDbStore

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
	"strings"
)

func LinkedListPush(
	db IGetPutLevelDb,
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
		newTail := []byte(head+"|"+strconv.Itoa(len(data)-1))
		err := db.Put(
			headKey,
			newTail,
			nil,
		)
		if err != nil {
			return err
		}
		err = db.Put(
			newTail,
			append([]byte(data[0]), []byte("(-1")...),
			nil,
		)
		for t:=1; t<len(data); t++ {
			err = db.Put(
				newTail,
				append(
					[]byte(data[0]),
					[]byte("(" + strconv.Itoa(t))...
				),
				nil,
			)
		}
		return err
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
		head+"|"+strconv.Itoa(index+len(data)-1),
	)
	err = db.Put(
		headKey,
		newTail,
		nil,
	)
	if err != nil {
		return err
	}
	for t:=index; t<index+len(data); t++ {
		err = db.Put(
			[]byte(head+"|"+strconv.Itoa(t)),
			append(
				[]byte(data[0]),
				[]byte("(" + strconv.Itoa(t))...
			),
			nil,
		)
	}

	return err
}
