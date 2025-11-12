package repository

import (
	"beta-be/internal/repository/ent"
	"beta-be/internal/repository/user"
)

type Registry interface {
	User() user.Repository
}

type impl struct {
	entConn *ent.Client
	userRepo user.Repository
}

func New(
	entConn *ent.Client,
) Registry {
	return impl{
		entConn: entConn,
		userRepo: user.New(entConn),
	}
}

func (i impl) User() user.Repository {
	return i.userRepo
}
