package user

import (
	"context"

	"beta-be/internal/model"
	"beta-be/internal/repository/ent"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Create(ctx context.Context, u model.User, passwordHash string) (model.User, error)
}

type impl struct {
	entClient *ent.Client
}

func New(entClient *ent.Client) Repository {
	return &impl{
		entClient: entClient,
	}
}
