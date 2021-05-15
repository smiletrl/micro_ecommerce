package logger

import (
	"context"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/kafka"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
)

func Consume(cfg config.KafkaConfig, logger logger.Provider, topic constants.KafkaTopic, partition int) error {
	kafka := kafka.NewProvider(cfg, logger)
	// @todo maybe define different consumer callbacks.
	err := kafka.Consume(context.Background(), topic, partition)
	return err
}
