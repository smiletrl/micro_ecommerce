package external

import (
	"context"
	"github.com/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	pb "github.com/smiletrl/micro_ecommerce/service.product/internal/rpc/proto"
	"google.golang.org/grpc"
	"time"
)

type Client interface {
	// Get sku detail
	GetSKU(skuID int64) (sku entity.SKU, err error)
}

type client struct {
	grpc pb.ProductClient
}

func NewClient() Client {
	// @todo use connection pool
	return client{}
}

// @todo add the connection pool
func newConnection() pb.ProductClient {
	// @todo inject this endpoint into config
	var productEndpoint = "product"
	var address = productEndpoint + constants.GrpcPort

	// need heart beat for this connection
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	//defer conn.Close()
	return pb.NewProductClient(conn)
}

func (c client) GetSKU(skuID int64) (sku entity.SKU, err error) {
	c.grpc = newConnection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.grpc.GetSKU(ctx, &pb.SKUID{Value: skuID})
	if err != nil {
		return sku, errors.Wrapf(errorsd.New("error getting sku from rpc"), "error getting sku from rpc: %s", err.Error())
	}
	if res != nil {
		sku = entity.SKU{
			Stock:      int(res.Stock),
			Amount:     int(res.Amount),
			Title:      res.Title,
			Attributes: res.Attributes,
		}
	}
	return sku, nil
}
