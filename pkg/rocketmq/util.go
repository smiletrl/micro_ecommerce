package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
)

// CreateMessage creates a single message.
// @todo maybe remove this function because its returned variable is allocated to heap.
func CreateMessage(topic constants.RocketMQTopic, tag constants.RocketMQTag, body string) *primitive.Message {
	message := primitive.NewMessage(string(topic), []byte(body))
	message.WithTag(string(tag))
	return message
}
