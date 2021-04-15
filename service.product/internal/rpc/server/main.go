package server

import (
	"context"
	"log"
	"net"

	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	pb "github.com/smiletrl/micro_ecommerce/service.product/internal/rpc/proto"
	"google.golang.org/grpc"
)

// Register the rpc server for product service.
func Register() {
	lis, err := net.Listen("tcp", constants.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// server is rpc server for product
type server struct {
	pb.UnimplementedProductServer
}

func (s *server) GetSKU(ctx context.Context, id *pb.SKUID) (*pb.SKU, error) {
	log.Printf("sku id is: %d\n", id.Value)

	// query db table directly
	// @todo get sku from cache firstly and then db.
	return &pb.SKU{
		Stock:      12,
		Amount:     1300,
		Title:      "Pretty desktop",
		Attributes: "size: XL; color: red",
	}, nil
}
