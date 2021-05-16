package balance

import (
	"context"
	"fmt"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

type Message interface {
	Subscribe() error
}

type message struct {
	// use map when multiple consumers available
	consumer mq.PushConsumer
	repo     Repository
	optMap   map[constants.RocketMQTag]messageOpt
	rocketmq rocketmq.Provider
	tracing  tracing.Provider
	logger   logger.Provider
}

type consumeOpt func(ctx context.Context, customerID int64, amount int) error

type messageOpt struct {
	consumeOpt  consumeOpt
	messageType entity.RocketmqMessage
}

func NewMessage(consumer mq.PushConsumer, repo Repository, rocketmq rocketmq.Provider, tracing tracing.Provider, logger logger.Provider) Message {
	optMap := map[constants.RocketMQTag]messageOpt{
		constants.RocketMQTagBalanceIncrease: messageOpt{
			consumeOpt:  repo.Increase,
			messageType: entity.RocketMQTagBalanceMessage{},
		},
		constants.RocketMQTagBalanceDecrease: messageOpt{
			consumeOpt:  repo.Decrease,
			messageType: entity.RocketMQTagBalanceMessage{},
		},
	}
	return message{consumer, repo, optMap, rocketmq, tracing, logger}
}

func (m message) Subscribe() error {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			m.logger.Errorw("rocketmq subscribe", err.Error())
		}
	}()

	err := m.subscribeDecreaseEvent()
	if err != nil {
		return err
	}

	err = m.subscribeIncreaseEvent()
	if err != nil {
		return err
	}
	return nil
}

func (m message) subscribeDecreaseEvent() error {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTagBalanceDecrease),
	}
	err := m.consumer.Subscribe(constants.RocketMQTopic, selector, m.callback(constants.RocketMQTagBalanceDecrease))
	return err
}

func (m message) subscribeIncreaseEvent() error {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTagBalanceIncrease),
	}
	err := m.consumer.Subscribe(constants.RocketMQTopic, selector, m.callback(constants.RocketMQTagBalanceIncrease))
	return err
}

func (m message) callback(tag constants.RocketMQTag) entity.RocketmqMessageOpt {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		var err error

		// Parse the message content and get the custom message type
		msg, ok := m.optMap[tag]
		if !ok {
			m.logger.Errorw("rocketmq balance message none-existing", string(tag))

			return consumer.Commit, err
		}
		rm, err := msg.messageType.Parse(string(msgs[0].Body))
		if err != nil {
			// This message should be sent to dead letter queue because the message self
			// is with incorrect format.
			m.logger.Errorw("rocketmq balance message", string(msgs[0].Body))

			return consumer.Commit, err
		}

		// See if this message has been consumed already.
		has, err := m.rocketmq.HasMessageConsumed(rm.Identifier())
		if err != nil {
			m.logger.Errorw("rocketmq balance message consumed", string(rm.Identifier()))

			return consumer.Commit, err
		}

		// If it has been consumed already, skip this message.
		if has {
			return consumer.ConsumeSuccess, nil
		}

		// Real consume happens here.
		err = msg.consumeOpt(ctx, rm.GetOption("customer_id").(int64), rm.GetOption("amount").(int))
		if err != nil {
			m.logger.Errorw("rocketmq balance opt invoke", err.Error())

			return consumer.ConsumeRetryLater, err
		}
		return consumer.ConsumeSuccess, nil
	}
}
