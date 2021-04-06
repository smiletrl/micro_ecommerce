package client

import (
	"context"
	"time"

	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	pb "github.com/smiletrl/micro_ecommerce/service.product/external/rpc"
	"google.golang.org/grpc"
)

const (
	// store this address in config
	address = "localhost:50051"
)

type Client interface {
	GetProduct(id int64) (prod entity.Product, err error)
}

type client struct {
	grpc pb.ProductClient
}

func NewClient() Client {
	// save this connection like db connection, so no need to build the connection all the time.
	// maybe save the connection in config, and then the connection will be built when parent service
	// is initialized.
	return client{}
}

func newConnection() pb.ProductClient {
	// need heart beat for this connection
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	//defer conn.Close()
	return pb.NewProductClient(conn)
}

func (c client) GetProduct(id int64) (prod entity.Product, err error) {
	// Build this connection somewhere. When this connection is broken, automatically
	// rebuild it.
	c.grpc = newConnection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.grpc.GetProduct(ctx, &pb.ProductID{Id: "12"})
	if err != nil {
		return prod, err
	}
	return entity.Product{
		Title:  res.Title,
		Stock:  12,
		Amount: 14,
	}, nil
}
