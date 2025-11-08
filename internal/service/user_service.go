package service

import (
	"context"

	"aiki/internal/domain"
	"aiki/internal/repository"
)

//go:generate mockgen -source=user_service.go -destination=mocks/mock_user_service.go -package=mocks

type UserService interface {
	GetByID(ctx context.Context, id int32) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, id int32, req *domain.UpdateUserRequest) (*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(ctx context.Context, id int32) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) Update(ctx context.Context, id int32, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.Update(ctx, id, req.FirstName, req.LastName, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}
