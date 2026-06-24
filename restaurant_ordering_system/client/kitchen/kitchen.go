package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"time"

	pb "github.com/iamads/go-workbook/restaurant_ordering_system/restaurant"
)

var serverAddress = "localhost:50051"

func printOrders(client pb.RestaurantClient) {
	log.Printf("Will stream in new orders as soon as they have been placed")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	stream, err := client.KitchenSubscribe(ctx, &emptypb.Empty{})

	if err != nil {
		log.Fatalf("client.KitchenSubscribe failed: %v", err)
	}

	for {
		order, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("client.Kitchen subscribe failed while receiving: %v", err)
		}
		log.Printf("New Order received: %s, for table %d", order.OrderItem.Name, order.TableNum)
	}
}

func main() {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRestaurantClient(conn)

	printOrders(client)
}
