package constants

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"strconv"
	"strings"
)

type RocketmqMessage interface {
	// used at producer
	SetOptions(options ...interface{}) RocketmqMessage
	String() string

	// used at consumer
	Parse(s string) (RocketmqMessage, error)
	GetOption(field string) interface{}
}

type MessageOpt func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)

type RocketMQGroup string

type RocketMQTag string

const (
	// RocketMQTopic rocketmq best practice is to only use one topic
	RocketMQTopic        string        = "micro_ecommerce"
	RocketMQGroupPayment RocketMQGroup = "payment"

	// Rocket tag
	// order
	RocketMQTagOrderPaid RocketMQTag = "order||paid"

	// balance
	RocketMQTagBalanceIncrease RocketMQTag = "balance||increase"
	RocketMQTagBalanceDecrease RocketMQTag = "balance||decrease"
)

// order
type RocketMQTagOrderPaidMessage struct {
	OrderID string
}

// options are: orderID
func (r RocketMQTagOrderPaidMessage) SetOptions(options ...interface{}) RocketmqMessage {
	rm := RocketMQTagOrderPaidMessage{
		OrderID: options[0].(string),
	}
	return rm
}

func (r RocketMQTagOrderPaidMessage) GetOption(field string) interface{} {
	if field == "order_id" {
		return r.OrderID
	}
	return nil
}

func (r RocketMQTagOrderPaidMessage) String() string {
	return fmt.Sprintf("order_id:%s", r.OrderID)
}

func (r RocketMQTagOrderPaidMessage) Parse(s string) (RocketmqMessage, error) {
	strSlice := strings.Split(s, ":")
	rm := RocketMQTagOrderPaidMessage{
		OrderID: strSlice[1],
	}
	return rm, nil
}

// balance
type RocketMQTagBalanceMessage struct {
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

func (r RocketMQTagBalanceMessage) String() string {
	return fmt.Sprintf("customer_id:%s||amount:%s", strconv.FormatInt(r.CustomerID, 10), strconv.Itoa(r.Amount))
}

func (r RocketMQTagBalanceMessage) Parse(s string) (RocketmqMessage, error) {
	rm := RocketMQTagBalanceMessage{}
	strSlice := strings.Split(s, "||")
	for _, str := range strSlice {
		strSubSlice := strings.Split(str, ":")
		if strSubSlice[0] == "customer_id" {
			customerID, err := strconv.ParseInt(strSubSlice[1], 10, 64)
			if err != nil {
				return rm, err
			}
			rm.CustomerID = customerID

		} else if strSubSlice[0] == "amount" {
			amount, err := strconv.Atoi(strSubSlice[1])
			if err != nil {
				return rm, err
			}
			rm.Amount = amount
		}
	}
	return rm, nil
}
