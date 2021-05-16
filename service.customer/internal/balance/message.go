package balance

import (
	"context"
	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

type Message interface {
	Subscribe() error
}

type message struct {
	// use map when multiple consumers available
	consumer   mq.PushConsumer
	repo       Repository
	optMap     map[constants.RocketMQTag]balanceOpt
	messageMap map[constants.RocketMQTag]constants.RocketmqMessage
	tracing    tracing.Provider
	logger     logger.Provider
}

type balanceOpt func(ctx context.Context, customerID int64, amount int) error

func NewMessage(consumer mq.PushConsumer, repo Repository, tracing tracing.Provider, logger logger.Provider) Message {
	optMap := map[constants.RocketMQTag]balanceOpt{
		constants.RocketMQTagBalanceIncrease: repo.Increase,
		constants.RocketMQTagBalanceDecrease: repo.Decrease,
	}

	msgMap := map[constants.RocketMQTag]constants.RocketmqMessage{
		constants.RocketMQTagBalanceIncrease: constants.RocketMQTagBalanceMessage{},
		constants.RocketMQTagBalanceDecrease: constants.RocketMQTagBalanceMessage{},
	}
	return message{consumer, repo, optMap, msgMap, tracing, logger}
}

func (m message) Subscribe() error {

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
	err := m.consumer.Subscribe(constants.RocketMQTopic, selector, m.opt(constants.RocketMQTagBalanceDecrease))
	return err
}

func (m message) subscribeIncreaseEvent() error {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTagBalanceIncrease),
	}
	err := m.consumer.Subscribe(constants.RocketMQTopic, selector, m.opt(constants.RocketMQTagBalanceIncrease))
	return err
}

func (m message) opt(tag constants.RocketMQTag) constants.MessageOpt {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		var err error

		// See if this message has been consumed already.

		msg, ok := m.messageMap[tag]
		if !ok {
			m.logger.Errorw("rocketmq balance message none-existing", string(tag))

			return consumer.Commit, err
		}
		rm, err := msg.Parse(string(msgs[0].Body))
		if err != nil {
			// This message should be sent to dead letter queue because the message self
			// is with incorrect format.
			m.logger.Errorw("rocketmq balance message", string(msgs[0].Body))

			return consumer.Commit, err
		}

		opt, ok := m.optMap[tag]
		if !ok {
			m.logger.Errorw("rocketmq balance opt none-existing", string(tag))

			return consumer.Commit, err
		}

		err = opt(ctx, rm.GetOption("customer_id").(int64), rm.GetOption("amount").(int))
		if err != nil {
			m.logger.Errorw("rocketmq balance opt invoke", err.Error())

			return consumer.ConsumeRetryLater, err
		}
		return consumer.ConsumeSuccess, nil
	}
}
