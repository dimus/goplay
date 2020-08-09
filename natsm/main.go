package main

import (
	"log"

	"google.golang.org/protobuf/proto"
)

const (
	NatsURL = "nats://pi:4222"
)

func main() {
	uri := NatsURL
	sd := &ServiceDiscovery{
		OrderServiceUri: uri,
	}
	psd, err := proto.Marshal(sd)
	if err != nil {
		log.Fatal(err)
	}
	_ = psd
}
