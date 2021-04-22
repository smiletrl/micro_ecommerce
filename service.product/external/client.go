package external

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	pb "github.com/smiletrl/micro_ecommerce/service.product/internal/rpc/proto"
	"google.golang.org/grpc"
	"time"
)

type Client interface {
	// Get sku stock
	GetSkuStock(eContext echo.Context, skuID string) (stock int, err error)

	// Get sku property
	GetSkuProperties(eContext echo.Context, skuIDs []string) (properties []entity.SkuProperty, err error)
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

func (c client) GetSkuStock(eContext echo.Context, skuID string) (stock int, err error) {
	c.grpc = newConnection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pbstock, err := c.grpc.GetSkuStock(ctx, &pb.SkuID{Value: skuID})
	if err != nil {
		return stock, errors.Wrapf(errorsd.New("error getting sku stock from rpc"), "error getting sku stock from rpc: %s", err.Error())
	}

	return int(pbstock.Value), nil
}

func (c client) GetSkuProperties(eContext echo.Context, skuIDs []string) (properties []entity.SkuProperty, err error) {
	c.grpc = newConnection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	gProperties, err := c.grpc.GetSkuProperties(ctx, &pb.SkuIDs{Value: skuIDs})
	if err != nil {
		return nil, errors.Wrapf(errorsd.New("error getting sku properties from rpc"), "error getting sku properties from rpc: %s", err.Error())
	}
	properties = make([]entity.SkuProperty, len(gProperties.Properties))
	for i, val := range gProperties.Properties {
		// maybe use int32 for entity
		properties[i] = entity.SkuProperty{
			SkuID:      val.GetId(),
			Title:      val.GetTitle(),
			Price:      int(val.GetPrice()),
			Attributes: val.GetAttributes(),
			Thumbnail:  val.GetThumbnail(),
			Stock:      int(val.GetStock()),
		}
	}

	return properties, nil
}
