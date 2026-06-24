package main

import (
	"context"
	"io"
	"log"
	"math/rand/v2"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"fmt"
	"time"

	pb "github.com/iamads/go-workbook/restaurant_ordering_system/restaurant"
)

var serverAddress = "localhost:50051"

// We get the menu print it
// and then return the menu
func printMenu(client pb.RestaurantClient) (*pb.Menu, error) {
	log.Println("Going to get menu")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	menu, err := client.GetMenu(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("client.GetMenu failed: %v", err)
		return nil, err
	}
	for _, item := range menu.Items {
		fmt.Printf("%s for Rs %d\n", item.Name, item.Price)
	}

	return menu, nil
}

// This initiates customers ordering flow
// Customer can stream their orders
// Once they are done they receive a bill
func customerOrder(client pb.RestaurantClient, menu *pb.Menu, tableNumber int) {
	log.Println("Pls place as much orders as you want!")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.PlaceOrder(ctx)
	if err != nil {
		log.Fatalf("client.PlaceOrder failed: %v", err)
	}
	myOrders := pickRandomOrders(menu, tableNumber) // A helper method I have added to pick random order

	for _, order := range myOrders {
		if err := stream.Send(order); err != nil {
			log.Fatalf("client.PlaceOrder failed at stream.send for order(%v) with error: %v", order, err)
			time.Sleep(500 * time.Millisecond)
		}
	}
	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.Place Order failed at stream close: %v", err)
	}
	log.Printf("Total bill for Table %d is %d \n", summary.TableNum, summary.Total)
}

// We use this function to review food
// This is a bidirectional rpc
// In this implementation customer sends a review and receives acknowlegment
// Note: Msg ordering is independent of client and server roles
func customerReview(client pb.RestaurantClient, tableNum int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.Review(ctx)
	if err != nil {
		log.Fatalf("client.Review failed: %v", err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("client.customerReview failed: %v", err)
			}
			log.Printf("Msg to customer: %s \n", in.Msg)
			if err = stream.CloseSend(); err != nil {
				log.Fatalf("Could not close stream after receving msg to customer: %v", err)
			}
		}
	}()

	review := getRandomReview()
	if err := stream.Send(&review); err != nil {
		log.Fatalf("client.Review: stream.Send(%v) failed: %v", review.Msg, err)
	}
	wg.Wait()
}

func getRandomReview() pb.ReviewChat {
	msgs := []string{"good!", "ok", "great!", "bad", "it was really bad"}
	selectedMsg := msgs[rand.IntN(len(msgs))]
	return pb.ReviewChat{Msg: selectedMsg}
}

func pickRandomOrders(menu *pb.Menu, tableNumber int) []*pb.Order {
	order := []*pb.Order{}

	// random number of items to be ordered

	itemCount := rand.IntN(7) + 1

	for i := 0; i < itemCount; i++ {
		order = append(order, &pb.Order{
			TableNum:  int32(tableNumber),
			OrderItem: menu.Items[rand.IntN(len(menu.Items))],
		})
	}
	return order
}

func main() {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRestaurantClient(conn)

	tableNumber := rand.IntN(10) + 1 // We have 10 tables and will be automatically assigned to customer

	menu, err := printMenu(client)
	if err != nil {
		log.Fatalf("Error in getting the menu: %v", err)
	}

	customerOrder(client, menu, tableNumber)

	log.Println("Thanks for dining at our restaurent! Pls share a review!")

	customerReview(client, tableNumber)

	log.Println("Shutting down client!")
}
