package payment

import (
	"context"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

type Message interface {
	ProduceOrderComplete(ctx context.Context, orderID string) error
	ProduceBalanceComplete(ctx context.Context, orderID string, customerID int64, amount int) error
}

type message struct {
	rocket  mq.Producer
	tracing tracing.Provider
}

func NewMessage(rocketP mq.Producer, tracing tracing.Provider) Message {
	return message{rocketP, tracing}
}

func (m message) ProduceOrderComplete(ctx context.Context, orderID string) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: ProduceOrderComplete order id:  "+orderID)
	defer m.tracing.FinishSpan(span)

	message := rocketmq.CreateMessage(constants.RocketMQTag("Pay Succeed||order"), "order_id:"+orderID)
	_, err := m.rocket.SendSync(ctx, message)
	return err
}

func (m message) ProduceBalanceComplete(ctx context.Context, orderID string, customerID int64, amount int) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: ProduceBalanceComplete customer id:  "+orderID)
	defer m.tracing.FinishSpan(span)

	message := rocketmq.CreateMessage(constants.RocketMQTag("Pay Succeed||balance"), "order_id:")
	_, err := m.rocket.SendSync(ctx, message)
	return err
}
