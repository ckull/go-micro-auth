package handler

import (
	"context"
	"go-meechok/config"
	"go-meechok/modules/inventory/model"
	"go-meechok/modules/inventory/useCase"
	"go-meechok/pkg/kafkaCon"
	"go-meechok/utils"
	"log"

	"github.com/segmentio/kafka-go"
)

type (
	InventoryQueueHandler interface {
		CreateInventory() error
	}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUsecase useCase.InventoryUsecase
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUsecase useCase.InventoryUsecase) InventoryQueueHandler {
	return &inventoryQueueHandler{
		cfg,
		inventoryUsecase,
	}
}

func (h *inventoryQueueHandler) Consumer(ctx context.Context, topic string) kafkaCon.Consumer {
	consumer := kafkaCon.NewConsumer(h.cfg.Kafka.Brokers, topic)

	return consumer
}

func (h *inventoryQueueHandler) Producer(ctx context.Context, topic string) kafkaCon.Producer {
	producer := kafkaCon.NewProducer(h.cfg.Kafka.Brokers, topic)

	return producer
}

func (h *inventoryQueueHandler) CreateInventory() error {
	ctx := context.Background()

	consumer := h.Consumer(ctx, kafkaCon.ProductTopic)

	defer consumer.Close()

	err := consumer.ConsumeMessages(ctx, func(msg kafka.Message) error {
		eventType := string(msg.Key)

		var req *model.AddInventoryReq
		switch eventType {
		case kafkaCon.ProductCreatedEvent:
			if err := utils.Deserialize(msg.Value, req); err != nil {
				log.Printf("Error deserializing message: %v", err)
				return err
			}

			if _, err := h.inventoryUsecase.AddInventory(req); err != nil {
				log.Printf("Error adding inventory: %v", err)

				failureEvent := &model.AddInventoryFailedEvent{
					ProductID: req.ProductID.Hex(),
				}
				producer := h.Producer(ctx, kafkaCon.InventoryTopic)
				defer producer.Close()

				failurePayload, err := utils.Serialize(failureEvent)
				if err != nil {
					return err
				}

				key, err := utils.Serialize(kafkaCon.InventoryCreatedFailedEvent)
				if err != nil {
					return err
				}

				if err := producer.Produce(ctx, key, failurePayload); err != nil {
					return err
				}

				return err

			}

			log.Println("Inventory successfully added for product:", req.ProductID)
			return nil
		default:
			log.Printf("default case")
			return nil
		}
	})

	if err != nil {
		return err
	}

	return nil

}
