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

func (s *server) GetSkuStock(ctx context.Context, skuID *pb.SkuID) (*pb.Stock, error) {
	log.Printf("sku id is: %s\n", skuID.Value)

	// query db table directly
	// @todo get sku from cache firstly and then db.
	return &pb.Stock{
		Value: int32(19),
	}, nil
}

func (s *server) GetSkuProperties(ctx context.Context, skuIDs *pb.SkuIDs) (*pb.SkuProperties, error) {
	log.Printf("sku ids are: %+v\n", skuIDs.Value)
	// @todo query the sku property from db.
	property := &pb.SkuProperty{
		Id:         "1233223",
		Stock:      12,
		Price:      12,
		Title:      "pretty desk",
		Attributes: "color: red, size: XL",
		Thumbnail:  "xx.png",
	}
	return &pb.SkuProperties{
		Properties: []*pb.SkuProperty{property},
	}, nil
}
