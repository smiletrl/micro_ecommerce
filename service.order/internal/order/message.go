package order

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
	// @todo add more produce message funcs
	//ProduceSkuDecrease(ctx context.Context, orderID string) error
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

type consumeOpt func(ctx context.Context, orderID string) error

type messageOpt struct {
	consumeOpt  consumeOpt
	messageType entity.RocketmqMessage
}

func NewMessage(consumer mq.PushConsumer, repo Repository, rocketmq rocketmq.Provider, tracing tracing.Provider, logger logger.Provider) Message {
	optMap := map[constants.RocketMQTag]messageOpt{
		constants.RocketMQTagOrderPaid: messageOpt{
			// @todo use service instead of repo. Need to be able to produce the message again.
			consumeOpt:  repo.OrderPaid,
			messageType: entity.RocketMQTagOrderPaidMessage{},
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

	err := m.subscribeOrderPaidEvent()
	if err != nil {
		return err
	}
	return nil
}

func (m message) subscribeOrderPaidEvent() error {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: string(constants.RocketMQTagOrderPaid),
	}
	err := m.consumer.Subscribe(constants.RocketMQTopic, selector, m.callback(constants.RocketMQTagOrderPaid))
	return err
}

func (m message) callback(tag constants.RocketMQTag) entity.RocketmqMessageOpt {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		var err error

		// Parse the message content and get the custom message type
		msg, ok := m.optMap[tag]
		if !ok {
			m.logger.Errorw("rocketmq order message none-existing", string(tag))

			return consumer.Commit, err
		}
		rm, err := msg.messageType.Parse(string(msgs[0].Body))
		if err != nil {
			// This message should be sent to dead letter queue because the message self
			// is with incorrect format.
			m.logger.Errorw("rocketmq order message", string(msgs[0].Body))

			return consumer.Commit, err
		}

		// See if this message has been consumed already.
		has, err := m.rocketmq.HasMessageConsumed(rm.Identifier())
		if err != nil {
			m.logger.Errorw("rocketmq order message consumed", string(rm.Identifier()))

			return consumer.Commit, err
		}

		// If it has been consumed already, skip this message.
		if has {
			return consumer.ConsumeSuccess, nil
		}

		// Real consume happens here.
		err = msg.consumeOpt(ctx, rm.GetOption("order_id").(string))
		if err != nil {
			m.logger.Errorw("rocketmq order opt invoke", err.Error())

			return consumer.ConsumeRetryLater, err
		}

		// Set the identifier consumed in db.
		if err := m.rocketmq.SetMessageConsumed(rm.Identifier()); err != nil {
			m.logger.Errorw("rocketmq balance identifier consumed", err.Error())

			// @todo need to workout a correct status
			return consumer.ConsumeSuccess, err
		}
		return consumer.ConsumeSuccess, nil
	}
}
