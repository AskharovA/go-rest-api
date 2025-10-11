package services

import (
	"AskharovA/go-rest-api/models"
	"AskharovA/go-rest-api/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.userRepo.Save(user)
}

func (s *UserService) ValidateCredentials(user *models.User) error {
	return s.userRepo.ValidateCredentials(user)
}
