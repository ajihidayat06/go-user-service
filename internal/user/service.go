package user

import (
	"context"
	"go-user-service/internal/pkg/errors"
)

type Service interface {
	Register(ctx context.Context, req CreateUserRequest) (*UserResponse, *errors.AppError)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, req CreateUserRequest) (*UserResponse, *errors.AppError) {
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // hash password di sini
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeDatabase, "Failed to create user")
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
