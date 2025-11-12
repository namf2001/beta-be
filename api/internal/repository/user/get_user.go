package user

import (
	"context"

	"beta-be/internal/model"
	"beta-be/internal/repository/ent"
	"beta-be/internal/repository/ent/user"

	pkgerrors "github.com/pkg/errors"
)

func (i impl) GetByEmail(ctx context.Context, email string) (model.User, error) {
	entUser, err := i.entClient.User.Query().Where(user.Email(email)).Only(ctx)
	if err != nil {
		return model.User{}, pkgerrors.WithStack(err)
	}
	return toModelUser(entUser), nil
}

func toModelUser(entUser *ent.User) model.User {
	return model.User{
		ID:       entUser.ID,
		Email:    entUser.Email,
		UserName: entUser.Username,
	}
}
