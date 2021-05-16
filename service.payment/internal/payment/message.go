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
	producer mq.Producer

	// rocketmq doesn't have jaeger natively supported yet
	// see https://github.com/apache/rocketmq/pull/1525
	tracing tracing.Provider
}

func NewMessage(producer mq.Producer, tracing tracing.Provider) Message {
	return message{producer, tracing}
}

func (m message) ProduceOrderComplete(ctx context.Context, orderID string) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: ProduceOrderComplete order id:  "+orderID)
	defer m.tracing.FinishSpan(span)

	message := rocketmq.CreateMessage(constants.RocketMQTag("Pay Succeed||order"), "order_id:"+orderID)
	_, err := m.producer.SendSync(ctx, message)
	return err
}

func (m message) ProduceBalanceComplete(ctx context.Context, orderID string, customerID int64, amount int) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: ProduceBalanceComplete customer id:  "+orderID)
	defer m.tracing.FinishSpan(span)

	message := rocketmq.CreateMessage(constants.RocketMQTag("Pay Succeed||balance"), "order_id:")
	_, err := m.producer.SendSync(ctx, message)
	return err
}
