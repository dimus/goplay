package main

import (
	"log"
	"time"

	nats "github.com/nats-io/nats"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func main() {

	// Create NATS server connection
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	log.Println("Connected to " + nats.DefaultURL)

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
		conn, err := grpc.Dial(address, grpc.WithInsecure())
	}
}
