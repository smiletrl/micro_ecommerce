package balance

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
	rocket := rocketmq.NewProvider(cfg)
	c, err := rocket.CreatePushConsumer(context.Background(), constants.RocketMQGroupPayment, consumer.Clustering)
	if err != nil {
		return err
	}
	selecter := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTag("Pay Succeed||method||customer||balance")),
	}
	err = c.Subscribe(string(constants.RocketMQTopicPayment), selecter, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		// Reduce the customer balance value
		fmt.Printf("subscribe payment callback in customer: %s \n", msgs[0].Body)
		return consumer.ConsumeSuccess, nil
	})
	return err
}
