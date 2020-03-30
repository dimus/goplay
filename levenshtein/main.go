package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dvirsky/levenshtein"
	"github.com/klauspost/compress/zstd"
	"gitlab.com/gogna/gnparser"
)

var decoder, _ = zstd.NewReader(nil)

func main() {
	trie := levenshtein.NewTrie()
	gnp := gnparser.NewGNparser()
	stems := getData()
	_ = stems
	buildTrie(trie, stems)
	raw := getRaw()
	log.Printf("start")
	count := 0
	for _, v := range raw {
		stem := getStem(gnp, v)
		if trie.Exists(v) {
			continue
		}
		count += 1
		matches := trie.FuzzyMatches(v, 2)
		fmt.Printf("\n%s: %s\n", v, stem)
		for _, vv := range matches {
			fmt.Printf("  %s: %s\n", v, vv)
		}
	}
	log.Printf("end")
	fmt.Println(count)
}

func getStem(gnp gnparser.GNparser, n string) string {
	p := gnp.ParseToObject(n)
	if p.Parsed {
		return p.Canonical.Stem
	} else {
		return ""
	}
}

func getRaw() []string {
	var names []string
	f, err := os.Open("data/raw-names.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	return names
}

func getData() []string {
	compressed, err := ioutil.ReadFile("data/stems.txt.zst")
	if err != nil {
		panic(err)
	}
	bs, err := decoder.DecodeAll(compressed, nil)
	if err != nil {
		panic(err)
	}
	bsr := bytes.NewReader(bs)
	scanner := bufio.NewScanner(bsr)
	canonicals := make([]string, 0, 10_000_000)
	for scanner.Scan() {
		canonicals = append(canonicals, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return canonicals
}

func buildTrie(trie *levenshtein.Trie, canonicals []string) {
	for _, v := range canonicals {
		trie.Insert(v)
	}
}
