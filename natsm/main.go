package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/dimus/goplay/natsm/pb"
	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
)

const (
	NatsURL = "nats://pi:4222"
)

func main() {
	conn0, err := nats.Connect(NatsURL)
	if err != nil {
		log.Fatal(err)
	}
	conn0.Subscribe("Discovery.OrderService", func(m *nats.Msg) {
		num := rand.Intn(100)
		val := strconv.Itoa(num)
		orderServiceDiscovery := pb.ServiceDiscovery{OrderServiceUri: val}
		data, err := proto.Marshal(&orderServiceDiscovery)
		fmt.Println(m.Respond)
		if err == nil {
			conn0.Publish(m.Reply, data)
		}
	})
	// Create NATS server connection
	natsConnection, err := nats.Connect(NatsURL)
	log.Println("Connected to " + NatsURL)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		msg, err := natsConnection.Request("Discovery.OrderService", nil, 1000*time.Millisecond)
		if err == nil && msg != nil {
			orderServiceDiscovery := pb.ServiceDiscovery{}
			err := proto.Unmarshal(msg.Data, &orderServiceDiscovery)
			if err != nil {
				log.Fatalf("Error on unmarshal: %v", err)
			}
			address := orderServiceDiscovery.OrderServiceUri
			log.Println("OrderService endpoint found at:", address)
			//Set up a connection to the gRPC server.
			// conn, err := grpc.Dial(address, grpc.WithInsecure())

		}
	}
}
