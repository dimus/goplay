package main

import (
	"fmt"
	"log"

	"github.com/gnames/gnfinder/protob"
	"github.com/golang/protobuf/proto"
)

func main() {
	v := protob.Pong{Value: "hi"}
	out, err := proto.Marshal(&v)
	if err != nil {
		log.Fatal(err)
	}
	res := protob.Pong{}
	err = proto.Unmarshal(out, &res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Value)

}
