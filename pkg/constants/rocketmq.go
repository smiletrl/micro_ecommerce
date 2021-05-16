package constants

type RocketMQGroup string

type RocketMQTag string

type MessageIdentifier string

const (
	// RocketMQTopic rocketmq best practice is to only use one topic
	RocketMQTopic        string        = "micro_ecommerce"
	RocketMQGroupPayment RocketMQGroup = "payment"

	RocketMQIdentifier string = "identifier"

	// Rocket tag
	// order
	RocketMQTagOrderPaid RocketMQTag = "order||paid"

	// balance
	RocketMQTagBalanceIncrease RocketMQTag = "balance||increase"
	RocketMQTagBalanceDecrease RocketMQTag = "balance||decrease"
)
