package main

import "fmt"

func main() {
	b := []byte("Homo")
	res, err := Parse("ddd", b)
	if err != nil {
		fmt.Println(err)
	}

	sn := res.(scientificNameNode)
	sn.Verbatim = string(b)
	fmt.Println(sn)
}

type Node interface {
	exec() error
}

type scientificNameNode struct {
	NamesGroupNode namesGroupNode
	Parsed         bool
	Verbatim       string
}

func newScientificNameNode(nameGr namesGroupNode) (scientificNameNode, error) {
	sn := scientificNameNode{NamesGroupNode: nameGr}
	return sn, nil
}

type namesGroupNode struct {
	NameNodes []nameNode
}

func newNamesGroupNode(names interface{}) (namesGroupNode, error) {
	items := toIfaceSlice(names)
	gn := make([]nameNode, len(items))
	for i, v := range items {
		gn[i] = v.(nameNode)
	}
	return namesGroupNode{NameNodes: gn}, nil
}

type nameNode struct {
	Uninomial uninomialNode
}

func newNameNode(n uninomialNode) (nameNode, error) {
	nn := nameNode{Uninomial: n}
	return nn, nil
}

type uninomialNode struct {
	Value string
}

func newUninomialNode(v string) (uninomialNode, error) {
	un := uninomialNode{Value: v}
	return un, nil
}
