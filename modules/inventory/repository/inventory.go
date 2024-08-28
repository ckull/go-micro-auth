package repository

import (
	"context"
	"fmt"
	"go-meechok/modules/inventory/model"
	"go-meechok/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryRepository interface {
		inventoryCollection() *mongo.Collection
		AddInventory(req *model.AddInventoryReq) (*model.Inventory, error)
		UpdateInventory(objectId primitive.ObjectID, amount int, version int) (*model.Inventory, error)
		RemoveInventory(inventoryId primitive.ObjectID) error
		FindOneInventoryByProductId(productId string) (*model.Inventory, error)
	}

	inventoryRepository struct {
		db    *mongo.Client
		redis *redis.Client
	}
)

func NewInventoryRepository(db *mongo.Client, redis *redis.Client) InventoryRepository {
	return &inventoryRepository{
		db,
		redis,
	}
}

func (r *inventoryRepository) inventoryCollection() *mongo.Collection {
	return r.db.Database("Inventory").Collection("Inventories")
}

func (r *inventoryRepository) AddInventory(req *model.AddInventoryReq) (*model.Inventory, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.inventoryCollection()

	newInventory := &model.Inventory{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.InsertOne(ctx, newInventory)
	if err != nil {
		return nil, err
	}

	return newInventory, err
}

func (r *inventoryRepository) AddInventoryWithTransaction(req *model.AddInventoryReq) (*model.Inventory, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newInventory *model.Inventory

	err := utils.Transaction(ctx, r.db, func(sc mongo.SessionContext) error {
		collection := r.inventoryCollection()

		newInventory := &model.Inventory{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := collection.InsertOne(ctx, newInventory)
		if err != nil {
			return err
		}

		return err
	})

	return newInventory, err
}

func (r *inventoryRepository) FindOneInventoryByProductId(productId string) (*model.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.inventoryCollection()

	filter := bson.M{"product_id": productId}
	inventory := new(model.Inventory)
	err := collection.FindOne(ctx, filter).Decode(&inventory)

	if err != nil {
		return nil, err
	}

	return inventory, err

}

func (r *inventoryRepository) UpdateInventory(objectId primitive.ObjectID, amount int, version int) (*model.Inventory, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.inventoryCollection()

	update := bson.M{
		"quantity": amount,
		"$inc": bson.M{
			"version": 1,
		}}
	filter := bson.M{
		"_id":     objectId,
		"version": version,
	}

	result := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedInventory model.Inventory
	if err := result.Decode(&updatedInventory); err != nil {
		return nil, err
	}

	return &updatedInventory, nil
}

func (r *inventoryRepository) RemoveInventory(inventoryId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.inventoryCollection()
	filter := bson.M{"_id": inventoryId}

	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no product found with the given ID")
	}

	return nil

}
