package payment

import (
	"context"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

type Message interface {
	ProduceOrderPaid(ctx context.Context, orderID string) error
	ProduceBalanceDecrease(ctx context.Context, orderID string, customerID int64, amount int) error
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

func (m message) ProduceOrderPaid(ctx context.Context, orderID string) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: produce order paid order id: "+orderID)
	defer m.tracing.FinishSpan(span)

	msg, err := rocketmq.NewMessage().Set("order_id", orderID).Encode(constants.RocketMQTagOrderPaid)
	if err != nil {
		return err
	}
	_, err = m.producer.SendSync(ctx, msg)
	return err
}

func (m message) ProduceBalanceDecrease(ctx context.Context, orderID string, customerID int64, amount int) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: ProduceBalanceComplete customer id: "+orderID)
	defer m.tracing.FinishSpan(span)

	msg, err := rocketmq.NewMessage().Set("order_id", orderID).Set("customer_id", customerID).Encode(constants.RocketMQTagBalanceDecrease)
	if err != nil {
		return err
	}

	_, err = m.producer.SendSync(ctx, msg)
	return err
}
