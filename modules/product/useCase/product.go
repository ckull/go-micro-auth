package useCase

import (
	"go-meechok/modules/product/model"
	"go-meechok/modules/product/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProductUsecase interface {
		CreateProduct(newProduct *model.Product) (*model.Product, error)
		CreateProductWithTransaction(req *model.CreateProduct) (*model.Product, error)
		RollbackProduct(productId primitive.ObjectID) error
	}

	productUsecase struct {
		productRepository repository.ProductRepository
	}
)

func NewProductUsecase(productRepository repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepository,
	}
}

func (u *productUsecase) CreateProduct(newProduct *model.Product) (*model.Product, error) {
	return u.productRepository.CreateProduct(newProduct)
}

func (u *productUsecase) CreateProductWithTransaction(req *model.CreateProduct) (*model.Product, error) {
	return u.productRepository.CreateProductWithTransaction(req)
}

func (u *productUsecase) RollbackProduct(productId primitive.ObjectID) error {
	return u.RollbackProduct(productId)
}
