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

func customerOrder(client pb.RestaurantClient, tableNumber int) {
	log.Println("Pls place as much orders as you want!")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.PlaceOrder(ctx)
	if err != nil {
		log.Fatalf("client.PlaceOrder failed: %v", err)
	}
	myOrders := pickRandomOrders() // A helper method I have added to pick random order

	myOrders = append(myOrders, &pb.Order{Done: true})
	for _, order := range myOrders {
		order.Id = int32(tableNumber)
		if err := stream.Send(order); err != nil {
			log.Fatalf("client.PlaceOrder failed at stream.send for order(%v) with error: %v", order, err)
			time.Sleep(500 * time.Millisecond)
		}
	}
	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.Place Order failed at stream close: %v", err)
	}
	log.Println("Your Complete Order: %v", summary)
}

func main() {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRestaurantClient(conn)

	printMenu(client)

	tableNumber := 4
	customerOrder(client, tableNumber)

	log.Println("Shutting down client!")
}

func pickRandomOrders() []*pb.Order {
	order := []*pb.Order{
		&pb.Order{
			Id:        0,
			OrderItem: &pb.MenuItem{Name: "Idli", Price: 20},
		},
	}
	return order
}
