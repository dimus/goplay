package main

import "testing"

func TestGetData(t *testing.T) {
	can := getData()
	if len(can) < 8_000_000 {
		t.Error("getData() should return an array of canonicals " +
			"larger than 8 million entries.")
	}
	if len(can) > 12_000_000 {
		t.Error("getData() should return an array of canonicals " +
			"it is smaller than 10 million entries.")
	}
}
