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
	_ "sync"
	"time"
)

type Provider interface {
	CreateProducer(ctx context.Context, group constants.RocketMQGroup) (rocketmq.Producer, error)
	CreatePushConsumer(ctx context.Context, group constants.RocketMQGroup, model consumer.MessageModel) (rocketmq.PushConsumer, error)

	HasMessageConsumed(id constants.MessageIdentifier) (bool, error)
	SetMessageConsumed(id constants.MessageIdentifier) error

	ShutdownProducer(producer rocketmq.Producer) error
	ShutdownPushConsumer(consumer rocketmq.PushConsumer) error
}

func NewProvider(cfg config.RocketMQConfig, pdb postgresql.Provider) Provider {
	p := provider{
		serverAddress: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.NameServerPort)},
		brokerAddress: fmt.Sprintf("%s:%s", cfg.Host, cfg.BrokerPort),
		pdb:           pdb,
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
	err = c.Start()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (p provider) HasMessageConsumed(id constants.MessageIdentifier) (bool, error) {
	// query p.pdb to get whether this identifier has been consumed from postgres
	return true, nil
}

func (p provider) SetMessageConsumed(id constants.MessageIdentifier) error {
	// query p.pdb to set identifier that it has been consumed
	return nil
}

func (p provider) ShutdownProducer(producer rocketmq.Producer) error {
	return producer.Shutdown()
}

func (p provider) ShutdownPushConsumer(consumer rocketmq.PushConsumer) error {
	return consumer.Shutdown()
}
