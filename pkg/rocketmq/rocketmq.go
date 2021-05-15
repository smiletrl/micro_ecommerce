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
)

type Provider interface {
	CreateDefaultTopic(ctx context.Context) error
	CreateProducer(ctx context.Context, group constants.RocketMQGroup) (rocketmq.Producer, error)
	CreatePushConsumer(ctx context.Context, group constants.RocketMQGroup, model consumer.MessageModel) (rocketmq.PushConsumer, error)
}

func NewProvider(cfg config.RocketMQConfig) Provider {
	return provider{
		serverAddress: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.NameServerPort)},
		brokerAddress: fmt.Sprintf("%s:%s", cfg.Host, cfg.BrokerPort),
	}
}

type provider struct {
	serverAddress []string
	brokerAddress string
}

func (p provider) CreateDefaultTopic(ctx context.Context) error {
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

/*
func Start(cfg config.RocketMQConfig) {
	var err error

	s := service{
		serverAddress: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.NameServerPort)},
		brokerAddress: fmt.Sprintf("%s:%s", cfg.Host, cfg.BrokerPort),
	}

	fmt.Printf("rocketmq service: %+v\n", s)

	// topic
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(p.serverAddress)))
	if err != nil {
		panic(err)
	}
	err = testAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate("jack"),
		admin.WithBrokerAddrCreate(p.brokerAddress),
	)
	if err != nil {
		panic(err)
	}
	// producer
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(p.serverAddress)),
		producer.WithRetry(2),
		producer.WithGroupName("GID_test"),
	)
	if err != nil {
		panic(err)
	}

	err = p.Start()
	if err != nil {
		panic(err)
	}
	message := primitive.NewMessage("jack", []byte("Hello Jack Go Client!"))
	message.WithTag("toml")

	_, err = p.SendSync(context.Background(), message)
	if err != nil {
		panic(err)
	}

	// consumer
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver(p.serverAddress)),
		consumer.WithGroupName("GID_test"),
		// model needs to be set after group name somehow to make topic filter working.
		consumer.WithConsumerModel(consumer.Clustering),
		//consumer.WithConsumerModel(consumer.BroadCasting),
	)

	err = c.Start()
	err = c.Subscribe("jack", consumer.MessageSelector{Type: consumer.TAG, Expression: "toml"}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe admin jack tom callback: %s \n", msgs[0].Body)
		return consumer.ConsumeSuccess, nil
	})
}
*/
