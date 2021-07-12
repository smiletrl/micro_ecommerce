package order

import (
	"context"
	"fmt"

	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/smiletrl/micro_ecommerce/pkg/constants"
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
	optMap   map[constants.RocketMQTag]consumeOpt
	rocketmq rocketmq.Provider
	tracing  tracing.Provider
	logger   logger.Provider
}

type consumeOpt func(ctx context.Context, orderID string) error

func NewMessage(consumer mq.PushConsumer, repo Repository, rocketmq rocketmq.Provider, tracing tracing.Provider, logger logger.Provider) Message {
	optMap := map[constants.RocketMQTag]consumeOpt{
		constants.RocketMQTagOrderPaid: repo.OrderPaid,
	}
	return message{consumer, repo, optMap, rocketmq, tracing, logger}
}

func (m message) Subscribe() error {
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

func (m message) callback(tag constants.RocketMQTag) rocketmq.MessageOpt {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				m.logger.Errorw("rocketmq subscribe callback error", err.Error())
			}
		}()

		rocketmsg, err := rocketmq.DecodeMessage(msgs[0].Body)
		if err != nil {
			return consumer.Commit, err
		}

		// See if this message has been consumed already.
		has, err := m.rocketmq.HasMessageConsumed(rocketmsg.ID())
		if err != nil {
			m.logger.Errorw("rocketmq order message consumed", string(rocketmsg.ID()))

			return consumer.Commit, err
		}

		// If it has been consumed already, skip this message.
		if has {
			return consumer.ConsumeSuccess, nil
		}

		// Get consume func
		consume, ok := m.optMap[tag]
		if !ok {
			m.logger.Errorw("rocketmq order message none-existing", string(tag))

			return consumer.Commit, err
		}

		// Real consume happens here.
		err = consume(ctx, rocketmsg.Get("order_id").(string))
		if err != nil {
			m.logger.Errorw("rocketmq order opt invoke", err.Error())

			return consumer.ConsumeRetryLater, err
		}

		// Set the message consumed in db.
		if err := m.rocketmq.SetMessageConsumed(rocketmsg.ID()); err != nil {
			m.logger.Errorw("rocketmq balance identifier consumed", err.Error())

			// @todo need to workout a correct status
			return consumer.ConsumeSuccess, err
		}
		return consumer.ConsumeSuccess, nil
	}
}
