package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	pb "github.com/smiletrl/micro_ecommerce/service.product/internal/rpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Register the rpc server for product service.
func Register(logger *zap.SugaredLogger) error {
	lis, err := net.Listen("tcp", constants.GrpcPort)
	if err != nil {
		logger.Errorf("error listening tcp connection for product grpc: %s", err.Error())
		return err
	}
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
	s := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			//grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_opentracing.UnaryServerInterceptor(),
		)),
	)
	pb.RegisterProductServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		logger.Errorf("error serving tcp connection for product grpc: %s", err.Error())
		return err
	}
	return nil
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
