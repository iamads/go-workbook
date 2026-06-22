package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"fmt"
	"time"

	pb "github.com/iamads/go-workbook/restaurant_ordering_system/restaurant"
)

var serverAddress = "localhost:50051"

func printMenu(client pb.RestaurantClient) {
	log.Println("Going to get menu")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	menu, err := client.GetMenu(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("client.GetMenu failed: %v", err)
	}
	for _, item := range menu.Items {
		fmt.Printf("%s for Rs %d\n", item.Name, item.Price)
	}
}

func main() {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRestaurantClient(conn)

	printMenu(client)

	log.Println("Shutting down client!")
}
