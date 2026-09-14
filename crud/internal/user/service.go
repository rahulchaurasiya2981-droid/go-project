package user

import "context"

type Service interface {
	GetUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, request CreateUserRequest) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) GetUsers(ctx context.Context) ([]User, error) {
	return s.repository.GetUsers(ctx)
}

func (s *service) CreateUser(ctx context.Context, request CreateUserRequest) (User, error) {
	return s.repository.CreateUser(ctx, request)
}
