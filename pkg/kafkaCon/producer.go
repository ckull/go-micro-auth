package kafkaCon

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

type (
	Producer interface {
		Produce(ctx context.Context, key, value []byte) error
		Close() error
	}
	producer struct {
		writer *kafka.Writer
	}
)

func NewProducer(brokers []string, topic string) Producer {
	return &producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
			Topic:    topic,
		}),
	}
}

func (p *producer) Produce(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

func (p *producer) Close() error {
	return p.writer.Close()
}

func (p *producer) DecodeMessage(obj any, value []byte) error {
	if err := json.Unmarshal(value, &obj); err != nil {
		log.Printf("Error: Failed to decode message: %s", err.Error())
		return errors.New("error: failed to decode message")
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		log.Printf("Error: Failed to validate message: %s", err.Error())
		return errors.New("error: failed to validate message")
	}

	return nil
}
