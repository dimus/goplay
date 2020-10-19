package binsrch

import (
	"fmt"
	"testing"
)

var (
	ary = []int{1, 3, 7, 10}
	tst = [][]int{
		{-3, -1},
		{1, 0},
		{6, -1},
		{7, 2},
		{10, 3},
		{106, -1},
	}
)

func TestBinsearch(t *testing.T) {
	for _, v := range tst {
		idx := binsearch(v[0], ary)
		if idx != v[1] {
			t.Errorf("%d should return idx %d, returned %d", v[0], v[1], idx)
		}
	}
}

func TestBinsearch2(t *testing.T) {
	for _, v := range tst {
		idx := binsearch2(v[0], ary)
		fmt.Printf("srch: %d\n", v[0])
		if idx != v[1] {
			t.Errorf("%d should return idx %d, returned %d", v[0], v[1], idx)
		}
	}
}
