package goplay

import (
	"fmt"
	"log"
)

type format int

const (
	Compact format = iota
	Pretty
	Simple
)

var formats = [...]string{"compact", "pretty", "simple"}

func (of format) String() string {
	return formats[of]
}

func newFormat(f string) format {
	for i, v := range formats {
		if v == f {
			return format(i)
		}
	}
	err := fmt.Errorf("Unknown format '%s', using default '%s' format.",
		f, Compact.String())
	log.Println(err)
	return Compact
}
