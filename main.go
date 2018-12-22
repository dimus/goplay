package main

import (
	"log"

	"github.com/dimus/goplay/grammar"
)

func main() {
	str := "Homo Linn 1999"
	gnp := &grammar.GNParser{Buffer: str, Pretty: true}
	gnp.Init()
	err := gnp.Parse()
	if err != nil {
		log.Println("ERR_START", err, "ERR_END")
	}
	gnp.PrintSyntaxTree()
}
