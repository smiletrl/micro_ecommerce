package product

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
)

func Consume(cfg config.RocketMQConfig) error {
	rocket := rocketmq.NewService(cfg)
	c, err := rocket.CreatePushConsumer(context.Background(), constants.RocketMQGroupPayment, consumer.Clustering)
	if err != nil {
		return err
	}
	selecter := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTag("Pay Succeed||product||sku||quantity")),
	}
	err = c.Subscribe(string(constants.RocketMQTopicPayment), selecter, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		// Reduce the product sku stock.
		fmt.Printf("subscribe payment callback in product: %s \n", msgs[0].Body)
		return consumer.ConsumeSuccess, nil
	})
	return err
}
