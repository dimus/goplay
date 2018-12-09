package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gnames/uuid5"
)

func main() {
	b := []byte("Pseudocercospora")
	res, err := Parse("", b)
	if err != nil {
		fmt.Println(err)
	}

	sn := res.(scientificNameNode)
	sn.addVerbatim(b)
	sn.addCanonical()
	output := newOutput(sn)
	json, err := output.toJSON(false)
	if err != nil {
		log.Fatal(err)
	}
	simple := newSimpleOutput(sn)
	ary := simple.toSlice()
	fmt.Println(string(json))
	fmt.Println(strings.Join(ary, "|"))

}

type Node interface {
	exec() error
}

type scientificNameNode struct {
	NamesGroup    []nameNode
	Parsed        bool
	Verbatim      string
	VerbatimID    string
	CanonicalNode canonicalNode
}

func newScientificNameNode(nameGr []nameNode) (scientificNameNode, error) {
	sn := scientificNameNode{NamesGroup: nameGr}
	return sn, nil
}

func (s *scientificNameNode) addVerbatim(bs []byte) {
	s.Verbatim = string(bs)
	s.VerbatimID = uuid5.UUID5(s.Verbatim).String()
}

func (s *scientificNameNode) addCanonical() {
	u := s.NamesGroup[0].Uninomial
	s.CanonicalNode = canonicalNode{
		Value:       u.Value,
		ValueRanked: u.Value,
	}
}

func newNamesGroup(names interface{}) ([]nameNode, error) {
	items := toIfaceSlice(names)
	gn := make([]nameNode, len(items))
	for i, v := range items {
		gn[i] = v.(nameNode)
	}
	return gn, nil
}

type nameNode struct {
	Uninomial uninomialNode
}

func newNameNode(n uninomialNode) (nameNode, error) {
	nn := nameNode{Uninomial: n}
	return nn, nil
}

type canonicalNode struct {
	Value       string
	ValueRanked string
}

type uninomialNode struct {
	Value string
}

func newUninomialNode(v string) (uninomialNode, error) {
	un := uninomialNode{Value: v}
	return un, nil
}
