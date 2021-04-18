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
)

type Service interface {
	CreateTopic(rocketMQ config.RocketMQConfig, topic string) error
	CreateProducer(rocketMQ config.RocketMQConfig, group string) error
	CreateConsumer(rocketMQ config.RocketMQConfig, topic, group string, model consumer.MessageModel) error
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
	fmt.Printf("move to consumer\n")
	// producer
	p, err := rocketmq.NewProducer(
		//producer.WithNameServer(endPoint),
		//producer.WithNsResovler(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		//producer.WithNsResovler(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
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
	result, err := p.SendSync(context.Background(), message)
	fmt.Printf("send sync result is: %+v\n", result)

	// consumer
	c, err := rocketmq.NewPushConsumer(
		//consumer.WithNameServer(endPoint),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		//consumer.WithConsumerModel(consumer.Clustering),
		//consumer.WithConsumerModel(consumer.BroadCasting),
		consumer.WithGroupName("GID_test"),
		// model needs to be set after group name somehow to make topic filter working.
		consumer.WithConsumerModel(consumer.Clustering),
	)

	err = c.Subscribe("jack", consumer.MessageSelector{Type: consumer.TAG, Expression: "toml"}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe jack tom callback: %s \n", msgs[0].Body)
		return consumer.ConsumeSuccess, nil
	})

	err = c.Start()
}
