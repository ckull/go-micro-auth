package repository

import (
	"context"
	"errors"
	"fmt"
	"go-meechok/modules/product/model"
	"go-meechok/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ProductRepository interface {
		productCollection() *mongo.Collection
		CreateProduct(req *model.Product) (*model.Product, error)
		removeProduct(productId primitive.ObjectID) error
		RollbackProduct(productId primitive.ObjectID) error
		CreateProductWithTransaction(req *model.CreateProduct) (*model.Product, error)
	}

	productRepository struct {
		db *mongo.Client
	}
)

func NewProductRepository(db *mongo.Client, redis *redis.Client) ProductRepository {
	return &productRepository{
		db,
	}
}

func (r *productRepository) productCollection() *mongo.Collection {
	return r.db.Database("Product").Collection("Products")
}

func (r *productRepository) CreateProduct(req *model.Product) (*model.Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newProduct *model.Product

	newProduct = &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Categories:  req.Categories,
		IsActive:    req.IsActive,
		SellerID:    req.SellerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := r.productCollection()

	_, err := collection.InsertOne(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	return newProduct, err

}

func (r *productRepository) RollbackProduct(productId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.productCollection()
	filter := bson.M{"_id": productId}
	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *productRepository) CreateProductWithTransaction(req *model.CreateProduct) (*model.Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newProduct *model.Product

	err := utils.Transaction(ctx, r.db, func(sc mongo.SessionContext) error {

		newProduct = &model.Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Categories:  req.Categories,
			IsActive:    req.IsActive,
			SellerID:    req.SellerID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		collection := r.productCollection()

		result, err := collection.InsertOne(ctx, newProduct)
		newProduct.ID = result.InsertedID.(primitive.ObjectID)
		return err

	})

	return newProduct, err

}

func (r *productRepository) removeProduct(productId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.productCollection()
	filter := bson.M{"_id": productId}

	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no product found with the given ID")
	}

	return nil

}
