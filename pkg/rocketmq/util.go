package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
)

// CreateMessage creates a single message.
// @todo maybe remove this function because its returned variable is allocated to heap.
func CreateMessage(tag constants.RocketMQTag, body string) *primitive.Message {
	message := primitive.NewMessage(constants.RocketMQTopic, []byte(body))
	message.WithTag(string(tag))
	return message
}
