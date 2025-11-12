package user

import (
	"beta-be/internal/repository"
)

type Controller interface {
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{
		repo: repo,
	}
}
