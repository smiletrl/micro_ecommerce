package client

import (
	"context"
	"log"
	"time"

	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	pb "github.com/smiletrl/micro_ecommerce/service.product/external/rpc"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetProduct(ctx, &pb.ProductID{Id: "12"})
	if err != nil {
		log.Fatalf("could not get product: %v", err)
	}
	log.Printf("Product title: %s", res.GetTitle())
}

func GetProduct(id int64) (entity.Product, error) {
	// save this connection like db connection, so no need to build the connection all the time.
	// maybe save the connection in config, and then the connection will be built when parent service
	// is initialized.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetProduct(ctx, &pb.ProductID{Id: "12"})
	return entity.Product{
		Title:  res.Title,
		Stock:  12,
		Amount: 14,
	}, nil
}
