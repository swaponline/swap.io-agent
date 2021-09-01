package journal

import (
	"testing"

	"swap.io-agent/src/blockchain"
)

func TestJournal(t *testing.T) {
	j := New("test")
	if len(j.GetSpends()) != 0 {
		t.Error(
			"invalid journal spends", j.GetSpends(),
			"expected - []",
		)
	}
	if len(j.GetSpendsAddress()) != 0 {
		t.Error(
			"invalid journal spends address", j.GetSpendsAddress(),
			"expected - []",
		)
	}

	j.Add("testId", blockchain.Spend{
		Wallet: "address",
		Value:  "1000",
		Label:  "label",
	})
	if len(j.GetSpends()) != 1 {
		t.Error(
			"invalid journal spends address", j.GetSpendsAddress(),
			"expected - []",
		)
	}
}
