package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/iamads/go-workbook/restaurant_ordering_system/restaurant"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port = 50051
)

type restaurentServer struct {
	pb.UnimplementedRestaurantServer

	mu       sync.Mutex
	orderMap map[int][]*pb.Order
}

func (s *restaurentServer) GetMenu(ctx context.Context, _ *emptypb.Empty) (*pb.Menu, error) {
	items := []*pb.MenuItem{
		&pb.MenuItem{Name: "Idli", Price: 20},
		&pb.MenuItem{Name: "Dosa", Price: 50},
		&pb.MenuItem{Name: "Biryani", Price: 200},
		&pb.MenuItem{Name: "Fried Rice", Price: 150},
	}
	menu := pb.Menu{
		Items: items,
	}
	return &menu, nil
}

func (s *restaurentServer) PlaceOrder(stream pb.Restaurant_PlaceOrderServer) error {
	for {
		order, err := stream.Recv()
		if err != nil {
			return err
		}

		// I am using done flag as to signify
		// customer is done ordering and would like
		// to see the bill
		//
		// TODO: ADD TOTAL Bill in summary
		if order.Done {
			allOrders := s.orderMap[int(order.Id)]

			orderItems := []*pb.MenuItem{}

			for _, item := range allOrders {
				orderItems = append(orderItems, item.OrderItem)
			}
			summary := pb.OrderSummary{
				OrderId:    order.Id,
				OrderItems: orderItems,
			}
			return stream.SendAndClose(&summary) // order is not there how will this happen
		}

		s.mu.Lock()
		if _, ok := s.orderMap[int(order.Id)]; ok {
			s.orderMap[int(order.Id)] = append(s.orderMap[int(order.Id)], order)
		} else {
			s.orderMap[int(order.Id)] = []*pb.Order{order}
		}
		s.mu.Unlock()
	}
}

func newServer() *restaurentServer {
	s := &restaurentServer{orderMap: make(map[int][]*pb.Order)}
	return s
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listed: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRestaurantServer(grpcServer, newServer())
	grpcServer.Serve(listener)
}
