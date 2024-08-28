package service

import (
	"go-meechok/config"
	"go-meechok/modules/inventory/model"
	"go-meechok/modules/product/useCase"
	"go-meechok/pkg/kafkaCon"
	"go-meechok/utils"
	"log"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProductQueueService interface {
		HandleInventoryCreatedFailedEvent(msg kafka.Message) error
		HandleInventoryEvent(msg kafka.Message) error
		NewConsumer(brokers string, topic string) kafkaCon.Consumer
	}

	productQueueService struct {
		productUsecase useCase.ProductUsecase
		cfg            *config.Config
	}
)

func NewProductQueueService(cfg *config.Config, productUsecase useCase.ProductUsecase) ProductQueueService {
	return &productQueueService{
		productUsecase,
		cfg,
	}
}

func (u *productQueueService) NewConsumer(brokers string, topic string) kafkaCon.Consumer {
	consumer := kafkaCon.NewConsumer(u.cfg.Kafka.Brokers, topic)
	return consumer
}

func (u *productQueueService) HandleInventoryEvent(msg kafka.Message) error {

	eventType := string(msg.Key)

	switch eventType {
	case kafkaCon.InventoryCreatedFailedEvent:
		return u.HandleInventoryCreatedFailedEvent(msg)

	}

	return nil

}

func (u *productQueueService) HandleInventoryCreatedFailedEvent(msg kafka.Message) error {

	var req *model.AddInventoryFailedEvent
	if err := utils.Deserialize(msg.Value, req); err != nil {
		log.Printf("Error deserializing message: %v", err)
		return err
	}

	objID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		log.Printf("Error convert objectID from Hex: %v", err)
		return err
	}

	if err := u.productUsecase.RollbackProduct(objID); err != nil {
		log.Printf("Error deserializing message: %v", err)
		return err
	}

	return nil
}
