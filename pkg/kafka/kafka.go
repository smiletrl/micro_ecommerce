package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"go.uber.org/zap"
	"time"
)

type Provider interface {
	CreateTopic(ctx context.Context, topic constants.KafkaTopic, partition int) error
	Produce(ctx context.Context, topic constants.KafkaTopic, partition int, messages ...string) error
	Consume(ctx context.Context, topic constants.KafkaTopic, partition int) error
}

func NewProvider(cfg config.KafkaConfig, logger *zap.SugaredLogger) Provider {
	return provider{
		url:    cfg.Host + cfg.Port,
		logger: logger,
	}
}

type provider struct {
	url    string
	logger *zap.SugaredLogger
}

func (p provider) CreateTopic(ctx context.Context, topic constants.KafkaTopic, partition int) error {
	// to create topics when auto.create.topics.enable='true'
	conn, err := kafka.DialLeader(ctx, "tcp", p.url, string(topic), partition)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (p provider) Produce(ctx context.Context, topic constants.KafkaTopic, partition int, messages ...string) error {
	// to produce messages
	conn, err := kafka.DialLeader(ctx, "tcp", p.url, string(topic), partition)
	if err != nil {
		p.logger.Errorf("error dailing kafka leader: %s", err.Error())
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	kMessages := make([]kafka.Message, len(messages))
	for i, message := range messages {
		kMessages[i] = kafka.Message{Value: []byte(message)}
	}
	_, err = conn.WriteMessages(kMessages...)
	if err != nil {
		p.logger.Errorf("error writing message in kafka: %s", err.Error())
		return nil
	}

	if err := conn.Close(); err != nil {
		p.logger.Errorf("error closing writer in kafka: %s", err.Error())
		return err
	}
	return nil
}

func (p provider) Consume(ctx context.Context, topic constants.KafkaTopic, partition int) error {
	conn, err := kafka.DialLeader(ctx, "tcp", p.url, string(topic), partition)
	if err != nil {
		p.logger.Errorf("error dailing kafka leader: %s", err.Error())
		return err
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			p.logger.Errorf("error reading message in kafka: %s", err.Error())
			break
		}
		fmt.Println(string(b))
	}

	if err := batch.Close(); err != nil {
		p.logger.Errorf("error closing batch in kafka: %s", err.Error())
		return err
	}

	if err := conn.Close(); err != nil {
		p.logger.Errorf("error closing connection in kafka: %s", err.Error())
		return err
	}
	return nil
}
