package external

import (
	"context"
	"github.com/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	pb "github.com/smiletrl/micro_ecommerce/service.product/internal/rpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Client interface {
	// Get sku stock
	GetSkuStock(ctx context.Context, skuID string) (stock int, err error)

	// Get sku property
	GetSkuProperties(ctx context.Context, skuIDs []string) (properties []entity.SkuProperty, err error)
}

type client struct {
	grpc   pb.ProductClient
	logger *zap.SugaredLogger
}

func NewClient(endpoint string, logger *zap.SugaredLogger) (Client, error) {
	conn, err := newConnectionClient(endpoint, logger)
	if err != nil {
		return nil, err
	}
	return client{
		grpc:   conn,
		logger: logger,
	}, nil
}

func newConnectionClient(endpoint string, logger *zap.SugaredLogger) (client pb.ProductClient, err error) {
	var address = endpoint + constants.GrpcPort

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	// only allow maximum 1 second connection.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		logger.Errorf("error connecting the grpc server at product: %s", err.Error())
		return nil, err
	}
	return pb.NewProductClient(conn), nil
}

func (c client) GetSkuStock(ctx context.Context, skuID string) (stock int, err error) {
	pbstock, err := c.grpc.GetSkuStock(ctx, &pb.SkuID{Value: skuID})
	if err != nil {
		return stock, errors.Wrapf(errorsd.New("error getting sku stock from rpc"), "error getting sku stock from rpc: %s", err.Error())
	}

	return int(pbstock.Value), nil
}

func (c client) GetSkuProperties(ctx context.Context, skuIDs []string) (properties []entity.SkuProperty, err error) {
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
