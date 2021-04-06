package server

import (
	"context"
	"log"
	"net"

	pb "github.com/smiletrl/micro_ecommerce/service.product/external/rpc"
	"google.golang.org/grpc"
)

const (
	// store the address somewhere else
	port = ":50051"
)

// server is rpc server for product
type server struct {
	pb.UnimplementedProductServer
}

func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.ProductDetail, error) {
	log.Printf("product id is: %s\n", in.Id)
	// get product from cache firstly and then db.
	// query db table directly
	return &pb.ProductDetail{Title: "mac", Amount: "1288", Stock: "12"}, nil
}

// build the server for product service.
// deploy to k8s.
func Register() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
