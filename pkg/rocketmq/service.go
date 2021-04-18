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

type Service interface {
	CreateTopic(topic constants.RocketMQTopic) error
	CreateProducer(group constants.RocketMQGroup) (rocketmq.Producer, error)
	CreatePushConsumer(group constants.RocketMQGroup, model consumer.MessageModel) (rocketmq.PushConsumer, error)
}

func NewService() Service {
	return service{}
}

type service struct {
	cfg config.RocketMQConfig
}

func (s service) CreateTopic(topic constants.RocketMQTopic) error {
	topicAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})))
	if err != nil {
		panic(err)
	}
	err = topicAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(string(topic)),
		admin.WithBrokerAddrCreate("127.0.0.1:10911"),
	)
	return err
}

func (s service) CreateProducer(group constants.RocketMQGroup) (rocketmq.Producer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
		producer.WithGroupName(string(group)),
	)
	if err != nil {
		return nil, err
	}
	err = p.Start()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s service) CreatePushConsumer(group constants.RocketMQGroup, model consumer.MessageModel) (rocketmq.PushConsumer, error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
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

func Start() {
	var err error

	// topic
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})))
	if err != nil {
		panic(err)
	}
	err = testAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate("jack"),
		admin.WithBrokerAddrCreate("127.0.0.1:10911"),
	)
	if err != nil {
		panic(err)
	}
	// producer
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
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
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
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
