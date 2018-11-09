package main

import (
	"github.com/dimus/goplay/util"
	"github.com/dimus/goplay/verifier"
)

func main() {
	names := []string{"Homo sapiens",
		"Pardosa moesta",
		"Felis concolor",
		"Plantago major alba"}
	m := util.NewModel()
	verifier.Verify(names, m)
}
