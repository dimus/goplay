package main

import (
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	NatsURL = "nats://pi:4222"
)

func main() {

	// Create NATS server connection
	natsConnection, _ := nats.Connect(NatsURL)
	log.Println("Connected to " + NatsURL)

	msg, err := natsConnection.Request("Discovery.OrderService", nil, 1000*time.Millisecond)
	if err == nil && msg != nil {
		orderServiceDiscovery := ServiceDiscovery{}
		err := proto.Unmarshal(msg.Data, &orderServiceDiscovery)
		if err != nil {
			log.Fatalf("Error on unmarshal: %v", err)
		}
		address := orderServiceDiscovery.OrderServiceUri
		log.Println("OrderService endpoint found at:", address)
		//Set up a connection to the gRPC server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
	}
}
