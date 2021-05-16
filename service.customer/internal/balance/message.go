package balance

import (
	"context"
	"fmt"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

type Message interface {
	Consume() error
}

type message struct {
	// use map when multiple consumers available
	consumer mq.PushConsumer
	tracing  tracing.Provider
	logger   logger.Provider
}

func NewMessage(consumer mq.PushConsumer, tracing tracing.Provider, logger logger.Provider) Message {
	return message{consumer, tracing, logger}
}

func (m message) Consume() error {
	err := m.consumeDecreaseEvent()
	if err != nil {
		return err
	}

	err = m.consumeIncreaseEvent()
	if err != nil {
		return err
	}
	return nil
}

func (m message) consumeDecreaseEvent() error {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTag("Pay Succeed||method||customer||balance")),
	}
	err := m.consumer.Subscribe(constants.RocketMQTopic, selector, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		// Reduce the customer balance value
		fmt.Printf("subscribe payment callback in customer: %s \n", msgs[0].Body)
		return consumer.ConsumeSuccess, nil
	})
	return err
}

func (m message) consumeIncreaseEvent() error {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTag("Pay Succeed||method||customer||balance")),
	}
	err := m.consumer.Subscribe(constants.RocketMQTopic, selector, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		// Reduce the customer balance value
		fmt.Printf("subscribe payment callback in customer: %s \n", msgs[0].Body)
		return consumer.ConsumeSuccess, nil
	})
	return err
}
