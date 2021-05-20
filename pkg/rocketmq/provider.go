package rocketmq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
	"time"
)

/*

Rocketmq provider is responsible for creating/shutdown rocket mq producer and consumer.

<b>Message idempotence</b>

It also checks whether one message has been consumed already to handle consumption idempotence. Right now, we
only use one postgres table to store the message ids. If this message id exists in this table, it means this
message has been consumed.

This message id will be primiary key in this table, maybe called `rocketmq_idempotence`. And the only column is `id`
to hold each message id.

It might be necessary to hold a date column `created_at` as well, so this table can be cleaned for message which has
been consumed 7 or 10 days ago.

This is actually an anti-pattern of micro service, because different services will make an update/query to the same table
`rocketmq_idempotence`. To make it follow micro service best practice, it might be necessary to create separate tables
such as `rocketmq_idempotence_order`, `rocketmq_idempotence_cart` to hold rocketmq message ids separately.

If each service is using its own postgres db, then it will connect to each service's own db to query table `rocketmq_idempotence`.

Depending on different needs, there might be other customized implementations.

<b>Topic</b>

According to [Rocketmq best practice](https://github.com/apache/rocketmq/blob/master/docs/cn/best_practice.md), only one topic
should be used for one app.

This provider is creating/initiating the default topic.
*/

type Provider interface {
	CreateProducer(ctx context.Context, group constants.RocketMQGroup) (rocketmq.Producer, error)
	CreatePushConsumer(ctx context.Context, group constants.RocketMQGroup, model consumer.MessageModel) (rocketmq.PushConsumer, error)

	HasMessageConsumed(id string) (bool, error)
	SetMessageConsumed(id string) error

	ShutdownProducer(producer rocketmq.Producer) error
	ShutdownPushConsumer(consumer rocketmq.PushConsumer) error
	StartPushConsumer(consumer rocketmq.PushConsumer) error
}

func NewProvider(cfg config.RocketMQConfig, pdb postgresql.Provider, serviceName string) Provider {
	p := provider{
		serverAddress: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.NameServerPort)},
		brokerAddress: fmt.Sprintf("%s:%s", cfg.Host, cfg.BrokerPort),
		pdb:           pdb,
		service:       serviceName,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// init the default topic
	err := p.createDefaultTopic(ctx)
	if err != nil {
		panic(err)
	}
	return p
}

type provider struct {
	serverAddress []string
	brokerAddress string
	service       string
	pdb           postgresql.Provider
}

func (p provider) createDefaultTopic(ctx context.Context) error {
	// @todo check if this topic existing already
	topicAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(p.serverAddress)))
	if err != nil {
		panic(err)
	}
	err = topicAdmin.CreateTopic(
		ctx,
		admin.WithTopicCreate(string(constants.RocketMQTopic)),
		admin.WithBrokerAddrCreate(p.brokerAddress),
	)
	return err
}

func (p provider) CreateProducer(ctx context.Context, group constants.RocketMQGroup) (rocketmq.Producer, error) {
	producer, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(p.serverAddress)),
		producer.WithRetry(2),
		producer.WithGroupName(string(group)),
	)
	if err != nil {
		return nil, err
	}
	err = producer.Start()
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (p provider) CreatePushConsumer(ctx context.Context, group constants.RocketMQGroup, model consumer.MessageModel) (rocketmq.PushConsumer, error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver(p.serverAddress)),
		consumer.WithGroupName(string(group)),
		// model needs to be set after group name somehow to make topic filter working.
		consumer.WithConsumerModel(model),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (p provider) HasMessageConsumed(id string) (bool, error) {
	// query p.pdb table `rocketmq_idempotence_{p.service}` to get whether this identifier has been consumed from postgres.
	return true, nil
}

func (p provider) SetMessageConsumed(id string) error {
	// query p.pdb table `rocketmq_idempotence_{p.service}` to set identifier that it has been consumed
	return nil
}

func (p provider) ShutdownProducer(producer rocketmq.Producer) error {
	return producer.Shutdown()
}

func (p provider) ShutdownPushConsumer(consumer rocketmq.PushConsumer) error {
	return consumer.Shutdown()
}
func (p provider) StartPushConsumer(consumer rocketmq.PushConsumer) error {
	return consumer.Start()
}
