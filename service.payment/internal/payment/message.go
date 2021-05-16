package payment

import (
	"context"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
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

	messageMap map[constants.RocketMQTag]entity.RocketmqMessage
}

func NewMessage(producer mq.Producer, tracing tracing.Provider) Message {
	// Init the rocketmq messsage map
	messageMap := map[constants.RocketMQTag]entity.RocketmqMessage{
		constants.RocketMQTagOrderPaid:       entity.RocketMQTagOrderPaidMessage{},
		constants.RocketMQTagBalanceDecrease: entity.RocketMQTagBalanceMessage{},
	}
	return message{producer, tracing, messageMap}
}

func (m message) ProduceOrderPaid(ctx context.Context, orderID string) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: ProduceOrderComplete order id: "+orderID)
	defer m.tracing.FinishSpan(span)

	msg := m.messageMap[constants.RocketMQTagOrderPaid]
	rm := msg.SetOptions(orderID)

	message := rocketmq.CreateMessage(constants.RocketMQTagOrderPaid, rm.String())
	_, err := m.producer.SendSync(ctx, message)
	return err
}

func (m message) ProduceBalanceDecrease(ctx context.Context, orderID string, customerID int64, amount int) error {
	span, ctx := m.tracing.StartSpan(ctx, "RocketMQ: ProduceBalanceComplete customer id: "+orderID)
	defer m.tracing.FinishSpan(span)

	msg := m.messageMap[constants.RocketMQTagBalanceDecrease]
	rm := msg.SetOptions(customerID, amount)

	message := rocketmq.CreateMessage(constants.RocketMQTagBalanceDecrease, rm.String())
	_, err := m.producer.SendSync(ctx, message)
	return err
}
