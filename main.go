package main

import "fmt"

type ScientificName struct {
	Parsed   bool
	Verbatim string
	Details  []item
}

type item struct {
	Uninomial uninomial
}

type uninomial struct {
	Value string
}

func main() {
	b := []byte("Homo")
	res, err := Parse("", b)
	if err != nil {
		fmt.Println(err)
	}

	sn := res.(ScientificName)
	sn.Verbatim = string(b)
	fmt.Println(sn)
}

func scientificName(n []byte) ScientificName {
	u := uninomial{Value: string(n)}
	det := []item{item{Uninomial: u}}
	res := ScientificName{Parsed: true, Details: det}
	return res
}
