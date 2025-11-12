package user

import (
	"context"
	"time"

	"beta-be/internal/model"

	pkgerrors "github.com/pkg/errors"
)

func (i impl) Create(ctx context.Context, u model.User, passwordHash string) (model.User, error) {
	created, err := i.entClient.User.Create().
		SetUsername(u.UserName).
		SetEmail(u.Email).
		SetPassword(passwordHash).
		SetActive(true).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return model.User{}, pkgerrors.WithStack(err)
	}
	return toModelUser(created), nil
}
