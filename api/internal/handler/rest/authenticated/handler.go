package authenticated

import (
	"beta-be/internal/controller/user"
)

type Handler struct {
	userCtrl user.Controller
}

// New instantiates a new Handler and returns it
func New(userCtrl user.Controller) Handler {
	return Handler{
		userCtrl: userCtrl,
	}
}
