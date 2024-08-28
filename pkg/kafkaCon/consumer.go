package kafkaCon

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type (
	Consumer interface {
		ConsumeMessages(ctx context.Context, handleMessage func(kafka.Message) error) error
		Close() error
		DecodeMessage(obj any, value []byte) error
	}

	consumer struct {
		reader *kafka.Reader
	}
)

func NewConsumer(brokers []string, topic string) Consumer {
	return &consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			// GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *consumer) ConsumeMessages(ctx context.Context, handleMessage func(kafka.Message) error) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Handle the message using the provided handler function
		err = handleMessage(msg)
		if err != nil {
			log.Printf("Failed to handle message: %v", err)
		}

		// Commit the message to Kafka
		c.reader.CommitMessages(ctx, msg)
		return err
	}
}
func (c *consumer) Close() error {
	return c.reader.Close()
}

func (c *consumer) DecodeMessage(obj any, value []byte) error {
	if err := json.Unmarshal(value, &obj); err != nil {
		log.Printf("Error: Failed to decode message: %s", err.Error())
		return errors.New("error: failed to decode message")
	}

	// validate := validator.New()
	// if err := validate.Struct(obj); err != nil {
	// 	log.Printf("Error: Failed to validate message: %s", err.Error())
	// 	return errors.New("error: failed to validate message")
	// }

	return nil
}
