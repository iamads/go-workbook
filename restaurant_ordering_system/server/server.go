package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "github.com/iamads/go-workbook/restaurant_ordering_system/restaurant"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port = 50051
	ch   = make(chan *pb.Order)
)

type restaurentServer struct {
	pb.UnimplementedRestaurantServer

	mu            sync.Mutex
	tableOrderMap map[int][]*pb.Order // record orders for the table
}

// return Menu to the client
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

// to place order for a cetain table
func (s *restaurentServer) PlaceOrder(stream pb.Restaurant_PlaceOrderServer) error {
	var tableNum int32
	for {
		order, err := stream.Recv()

		if err == io.EOF {
			allOrders := s.tableOrderMap[int(tableNum)]

			orderItems := []*pb.MenuItem{}
			total := 0

			for _, item := range allOrders {
				orderItems = append(orderItems, item.OrderItem)
				total += int(item.OrderItem.Price)
			}
			summary := pb.OrderSummary{
				TableNum:   tableNum,
				OrderItems: orderItems,
				Total:      int32(total),
			}
			return stream.SendAndClose(&summary)
		}

		if err != nil {
			return err
		}

		if tableNum == 0 {
			tableNum = order.TableNum
		}

		s.mu.Lock()
		if _, ok := s.tableOrderMap[int(order.TableNum)]; ok {
			s.tableOrderMap[int(order.TableNum)] = append(s.tableOrderMap[int(order.TableNum)], order)
		} else {
			s.tableOrderMap[int(order.TableNum)] = []*pb.Order{order}
		}
		s.mu.Unlock()
		ch <- order // send the order to kitchen
	}
}

// Newly placed order will be streamed to client
func (s *restaurentServer) KitchenSubscribe(_ *emptypb.Empty, stream pb.Restaurant_KitchenSubscribeServer) error {
	for order := range ch {
		if err := stream.Send(order); err != nil {
			return err
		}
	}

	return nil
}

// We create a bidirectional stream for collecting Review
// customer sends their review
// we thank them for response and clear out their table related info from orderMap
func (s *restaurentServer) Review(stream pb.Restaurant_ReviewServer) error {
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Println("Customer says: ", in.Msg)

		if err := stream.Send(&pb.ReviewChat{Msg: "Thanks for sharing your experience!"}); err != nil {
			return err
		} else {
			delete(s.tableOrderMap, int(in.TableNum))
		}
	}
}

func newServer() *restaurentServer {
	s := &restaurentServer{tableOrderMap: make(map[int][]*pb.Order)}
	return s
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listed: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRestaurantServer(grpcServer, newServer())

	log.Println("Starting server")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to start listening: %v", err)
	}
}
