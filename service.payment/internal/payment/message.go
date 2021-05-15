package payment

import (
	"context"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
)

type Message interface {
	SendOrderComplete(ctx context.Context, orderID string) error
	SendBalanceComplete(ctx context.Context, customerID int64, amount int) error
}

type message struct {
	rocket mq.Producer
}

func NewMessage(rocketP mq.Producer) Message {
	return message{rocketP}
}

func (m message) SendOrderComplete(ctx context.Context, orderID string) error {
	message := rocketmq.CreateMessage(constants.RocketMQTag("Pay Succeed||order"), "order_id:"+orderID)
	_, err := m.rocket.SendSync(ctx, message)
	return err
}

func (m message) SendBalanceComplete(ctx context.Context, customerID int64, amount int) error {
	message := rocketmq.CreateMessage(constants.RocketMQTag("Pay Succeed||balance"), "order_id:")
	_, err := m.rocket.SendSync(ctx, message)
	return err
}
