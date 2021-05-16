package entity

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"strconv"

	"strings"
)

type RocketmqMessage interface {
	// Used at producer
	SetOptions(options ...interface{}) RocketmqMessage
	// String should be called after SetOptions, and it is used to be sent to queue.
	String() string

	// Used at consumer
	// Parse should be called before GetOption and HasConsumed.
	Parse(s string) (RocketmqMessage, error)
	GetOption(field string) interface{}
	Identifier() constants.MessageIdentifier
}

type RocketmqMessageOpt func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)

// util
func MessageUUID(s string) string {
	return fmt.Sprintf("%s:%s||%s", constants.RocketMQIdentifier, uuid.New().String(), s)
}

// order
type RocketMQTagOrderPaidMessage struct {
	constants.MessageIdentifier
	OrderID string
}

// options are: orderID
func (r RocketMQTagOrderPaidMessage) SetOptions(options ...interface{}) RocketmqMessage {
	rm := RocketMQTagOrderPaidMessage{
		OrderID: options[0].(string),
	}
	return rm
}

func (r RocketMQTagOrderPaidMessage) String() string {
	return MessageUUID(fmt.Sprintf("order_id:%s", r.OrderID))
}

func (r RocketMQTagOrderPaidMessage) Parse(s string) (RocketmqMessage, error) {
	rm := RocketMQTagOrderPaidMessage{}
	strSlice := strings.Split(s, "||")

	for _, str := range strSlice {
		// @todo uuid might include ":" as well.
		strSubSlice := strings.Split(str, ":")
		switch strSlice[0] {
		case constants.RocketMQIdentifier:
			rm.MessageIdentifier = constants.MessageIdentifier(strSlice[1])
		case "order_id":
			rm.OrderID = strSubSlice[1]
		}
	}
	return rm, nil
}

func (r RocketMQTagOrderPaidMessage) GetOption(field string) interface{} {
	switch field {
	case "order_id":
		return r.OrderID
	default:
		return nil
	}
}

func (r RocketMQTagOrderPaidMessage) Identifier() constants.MessageIdentifier {
	return r.MessageIdentifier
}

// balance
type RocketMQTagBalanceMessage struct {
	constants.MessageIdentifier
	CustomerID int64
	Amount     int
}

// options are: customerID, amount
func (r RocketMQTagBalanceMessage) SetOptions(options ...interface{}) RocketmqMessage {
	rm := RocketMQTagBalanceMessage{
		CustomerID: options[0].(int64),
		Amount:     options[1].(int),
	}
	return rm
}

func (r RocketMQTagBalanceMessage) String() string {
	return MessageUUID(fmt.Sprintf("customer_id:%s||amount:%s", strconv.FormatInt(r.CustomerID, 10), strconv.Itoa(r.Amount)))
}

func (r RocketMQTagBalanceMessage) Parse(s string) (RocketmqMessage, error) {
	rm := RocketMQTagBalanceMessage{}
	strSlice := strings.Split(s, "||")

	for _, str := range strSlice {
		strSubSlice := strings.Split(str, ":")

		switch strSlice[0] {
		case constants.RocketMQIdentifier:
			rm.MessageIdentifier = constants.MessageIdentifier(strSlice[1])
		case "customer_id":
			customerID, err := strconv.ParseInt(strSubSlice[1], 10, 64)
			if err != nil {
				return rm, err
			}
			rm.CustomerID = customerID
		case "amount":
			amount, err := strconv.Atoi(strSubSlice[1])
			if err != nil {
				return rm, err
			}
			rm.Amount = amount
		}
	}
	return rm, nil
}

func (r RocketMQTagBalanceMessage) GetOption(field string) interface{} {
	switch field {
	case "customer_id":
		return r.CustomerID
	case "amount":
		return r.Amount
	default:
		return nil
	}
}

func (r RocketMQTagBalanceMessage) Identifier() constants.MessageIdentifier {
	return r.MessageIdentifier
}
