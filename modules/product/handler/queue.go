package handler

import (
	"context"
	"go-meechok/config"
	inventoryModel "go-meechok/modules/inventory/model"
	"go-meechok/modules/product/model"
	productModel "go-meechok/modules/product/model"
	"go-meechok/modules/product/service"
	"go-meechok/modules/product/useCase"
	"go-meechok/pkg/kafkaCon"
	"go-meechok/utils"
	"log"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProductQueueHandler interface {
		Producer(ctx context.Context, topic string) kafkaCon.Producer
		CreateProduct(req *model.CreateProduct) error
	}

	productQueueHandler struct {
		cfg                 *config.Config
		productUsecase      useCase.ProductUsecase
		productQueueService service.ProductQueueService
	}
)

func NewProductQueueHandler(
	cfg *config.Config,
	productUsecase useCase.ProductUsecase,
	productQueueService service.ProductQueueService,
) ProductQueueHandler {
	return &productQueueHandler{
		cfg,
		productUsecase,
		productQueueService,
	}
}

func (h *productQueueHandler) Producer(ctx context.Context, topic string) kafkaCon.Producer {
	producer := kafkaCon.NewProducer(h.cfg.Kafka.Brokers, topic)

	return producer
}

func (h *productQueueHandler) CreateProduct(req *productModel.CreateProduct) error {
	ctx := context.Background()

	newProduct, err := h.productUsecase.CreateProductWithTransaction(req)
	if err != nil {
		return err
	}

	producer := h.Producer(ctx, kafkaCon.ProductTopic)

	defer producer.Close()

	productEvent := &productModel.ProductCreatedPayload{
		ProductID: newProduct.ID.Hex(),
		Quantity:  req.Quantity,
	}

	payload, err := utils.Serialize(productEvent)
	if err != nil {
		return err
	}

	key, err := utils.Serialize(kafkaCon.ProductCreatedEvent)
	if err != nil {
		return err
	}

	if err := producer.Produce(ctx, key, payload); err != nil {
		return err
	}

	go h.HandleInventoryResponse(ctx)

	return err

}

func (h *productQueueHandler) HandleInventoryResponse(ctx context.Context) {
	consumer := kafkaCon.NewConsumer(h.cfg.Kafka.Brokers, kafkaCon.InventoryTopic)
	defer consumer.Close()

	var req *inventoryModel.AddInventoryFailedEvent
	handleMessage := func(msg kafka.Message) error {

		eventType := string(msg.Key)

		switch eventType {
		case kafkaCon.InventoryCreatedFailedEvent:
			if err := utils.Deserialize(msg.Value, req); err != nil {
				log.Printf("Error deserializing message: %v", err)
				return err
			}

			objID, err := primitive.ObjectIDFromHex(req.ProductID)
			if err != nil {
				log.Printf("Error convert objectID from Hex: %v", err)
				return err
			}
			if err := h.productUsecase.RollbackProduct(objID); err != nil {
				log.Printf("Error deserializing message: %v", err)
				return err
			}

		}

		return nil
	}

	err := consumer.ConsumeMessages(ctx, handleMessage)
	if err != nil {
		log.Printf("Error in handling inventory response: %v", err)
		// Depending on the error, you might want to break the loop or continue
		// For now, we'll continue and retry

	}

}
