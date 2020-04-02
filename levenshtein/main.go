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
	raw := getRaw(gnp, trie)
	log.Printf("Found %d stem to check", len(raw))
	log.Println("start")
	count := 0
	for k, v := range raw {
		count += 1
		if count%10_000 == 0 {
			log.Printf("Processing %d-th row", count)
		}
		matches := trie.FuzzyMatches(v, 2)
		fmt.Printf("\n%s:   %s\n", k, v)
		for _, vv := range matches {
			fmt.Printf("  %s: %s\n", k, vv)
		}
	}
	log.Println("end")
	fmt.Println(len(raw))
}

func getStem(gnp gnparser.GNparser, n string) string {
	p := gnp.ParseToObject(n)
	if p.Parsed {
		return p.Canonical.Stem
	} else {
		return ""
	}
}

func getRaw(gnp gnparser.GNparser, trie *levenshtein.Trie) map[string]string {
	names := make(map[string]string)
	f, err := os.Open("data/raw-names.txt")
	if err != nil {
		panic(err)
	}
	count := 0
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name := scanner.Text()
		stem := getStem(gnp, name)
		if stem != "" && !trie.Exists(stem) {
			names[name] = stem
		} else {
			count += 1
		}
	}
	log.Printf("Found %d exact stem matches\n", count)
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
	stems := make([]string, 0, 10_000_000)
	for scanner.Scan() {
		stems = append(stems, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return stems
}

func buildTrie(trie *levenshtein.Trie, stems []string) {
	for _, v := range stems {
		trie.Insert(v)
	}
}
