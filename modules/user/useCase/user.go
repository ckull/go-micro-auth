package useCase

import (
	"go-auth/modules/user/model"
	"go-auth/modules/user/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserUsecase interface {
		GetUserByUID(uid primitive.ObjectID) (model.User, error)
	}

	userUsecase struct {
		userRepository repository.UserRepository
	}
)

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (uc *userUsecase) GetUserByUID(uid primitive.ObjectID) (model.User, error) {
	return uc.userRepository.GetUserByUID(uid)
}
