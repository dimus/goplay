package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/klauspost/compress/zstd"
)

var decoder, _ = zstd.NewReader(nil)

func main() {
	data := getData()
	canonicals := strings.Split(string(data), "\n")
	for _, v := range canonicals {
		fmt.Printf("'%s'\n", v)
	}

}

func getData() []byte {
	canonicals, err := ioutil.ReadFile("data/canonicals.txt.zst")
	if err != nil {
		panic(err)
	}
	bs, err := decoder.DecodeAll(canonicals, nil)
	if err != nil {
		panic(err)
	}
	return bs
}
