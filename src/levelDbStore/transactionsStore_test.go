package levelDbStore

import (
	"os"
	"testing"
)

const storeNameTest = "testWrite"

func clearTestAffect() {
	os.RemoveAll(transactionsStoreDbDir)
}

func TestWriteTransactionStore(t *testing.T) {
	clearTestAffect()

	store, err := InitialiseTransactionStore(TransactionsStoreConfig{
		Name:                 storeNameTest,
		DefaultScannedBlocks: -1,
	})
	if err != nil {
		t.Error("not create transaction store", err)
	}
	if store.lastBlock != -1 {
		t.Error("last block not set default value")
	}

	store.WriteLastIndexedTransactions(map[string][]string{
		"add":  {"hash", "hash2"},
		"add1": {"hash1", "hash3"},
	}, 0)
	if store.lastBlock != 0 {
		t.Error("last block not updated")
	}
	err = store.Flush()
	if err != nil {
		t.Error("not write", err)
	}
	store.WriteLastIndexedTransactions(map[string][]string{
		"add":  {"hash4", "hash6"},
		"add1": {"hash5", "hash7"},
		"add3": {},
	}, 1)
	err = store.Flush()
	if err != nil {
		t.Error("not write", err)
	}

	notExistCursor, err := store.GetCursorFromAddress("add3")
	if err != nil {
		t.Error(err)
	}
	if notExistCursor != "null" {
		t.Error("cursor generated then 0 hashes write")
	}
	notExistCursor1, err := store.GetCursorFromAddress("notExist")
	if err != nil {
		t.Error(err)
	}
	if notExistCursor1 != "null" {
		t.Error("cursor generated then 0 hashes write")
	}
	notExistCursorData, err := store.GetCursorTransactionHashes("notExist")
	if err != nil {
		t.Error(err)
	}
	if notExistCursorData.NextCursor != "null" {
		t.Error("cursor generated then 0 hashes write")
	}
	if len(notExistCursorData.Hashes) != 0 {
		t.Error("not exist cursor data generated")
	}

	cursor, err := store.GetCursorFromAddress("add1")
	if err != nil {
		t.Error(err)
	}
	if cursor == "null" {
		t.Error("cursor not generated then 4 hashes write")
	}
	cursorData, err := store.GetCursorTransactionHashes(cursor)
	if err != nil {
		t.Error(err)
	}
	if cursorData.Cursor != "add1|1" {
		t.Error(
			"invalid cursor", cursorData.Cursor,
			"expected", "add1|1",
		)
	}
	if cursorData.NextCursor != "add1|0" {
		t.Error(
			"invalid next cursor", cursorData.NextCursor,
			"expected", "add1|0",
		)
	}
	if len(cursorData.Hashes) != 2 {
		t.Error("cursor data expected 2")
	}
	if cursorData.Hashes[0] != "hash5" || cursorData.Hashes[1] != "hash7" {
		t.Error(
			"incorrect cursor data, expected [hash5, hash7]",
			"got get ", cursorData.Hashes,
		)
	}
	cursorData1, err := store.GetCursorTransactionHashes(cursorData.NextCursor)
	if err != nil {
		t.Error(err)
	}
	if cursorData1.Cursor != "add1|0" {
		t.Error(
			"invalid cursor", cursorData1.Cursor,
			"expected", "add1|0",
		)
	}
	if cursorData1.NextCursor != "null" {
		t.Error(
			"invalid next cursor", cursorData1.NextCursor,
			"expected", "null",
		)
	}
	if len(cursorData1.Hashes) != 2 {
		t.Error("cursor data expected 2")
	}
	if cursorData1.Hashes[0] != "hash1" || cursorData1.Hashes[1] != "hash3" {
		t.Error(
			"incorrect cursor data, expected [hash1, hash3]",
			"got get ", cursorData1.Hashes,
		)
	}
	cursorData2, err := store.GetCursorTransactionHashes(cursorData1.NextCursor)
	if err != nil {
		t.Error(err)
	}
	if cursorData2.Cursor != "null" {
		t.Error(
			"invalid cursor", cursorData2.Cursor,
			"expected", "null",
		)
	}
	if cursorData2.NextCursor != "null" {
		t.Error(
			"invalid next cursor", cursorData2.NextCursor,
			"expected", "null",
		)
	}
	if len(cursorData2.Hashes) != 0 {
		t.Error("cursor data expected 0")
	}

	clearTestAffect()
}
