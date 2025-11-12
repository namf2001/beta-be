package user

import (
	"context"

	"beta-be/internal/model"
	"beta-be/internal/repository"
)

type Controller interface {
	Register(ctx context.Context, user model.User, passwordHash string) (model.User, error)
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{
		repo: repo,
	}
}

func (i impl) Register(ctx context.Context, user model.User, passwordHash string) (model.User, error) {
	// Create user in repository
	createdUser, err := i.repo.User().Create(ctx, user, passwordHash)
	if err != nil {
		return model.User{}, err
	}
	return createdUser, nil
}
