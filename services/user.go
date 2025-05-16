package services

import (
	"fmt"

	"app/models"
	"app/repositories"
	"app/utils"
)

type UserService interface {
	SignUp(username, password string) (*models.User, error)
	SignIn(username, password string) (string, *models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) SignUp(username, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &models.User{
		Username: username,
		Password: hashedPassword,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *userService) SignIn(username, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, fmt.Errorf("invalid username: %w", err)
	}
	if err := utils.VerifyPassword(user.Password, password); err != nil {
		return "", nil, fmt.Errorf("invalid password: %w", err)
	}
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}
	return token, user, nil
}
