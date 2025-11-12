package router

import (
	"context"
	"net/http"

	"beta-be/internal/controller/user"
	"beta-be/internal/handler/rest/jwt"

	jwtmiddleware "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Router struct {
	ctx           context.Context
	jwtHandler    jwt.Handler
	jwtMiddleware *jwtmiddleware.GinJWTMiddleware
}

func New(
	ctx context.Context,
	userCtrl user.Controller,
	jwtMiddleware *jwtmiddleware.GinJWTMiddleware,
) Router {
	r := Router{
		ctx:           ctx,
		jwtHandler:    jwt.New(userCtrl),
		jwtMiddleware: jwtMiddleware,
	}
	r.jwtMiddleware.Authenticator = r.jwtHandler.Authenticator
	return r
}
func (rtr Router) Handler(g *gin.Engine) (http.Handler, error) {
	rtr.public(g)
	rtr.jwt(g)
	return g, nil
}

func (rtr Router) public(g *gin.Engine) {

}

func (rtr Router) jwt(g *gin.Engine) {
	v1 := g.Group("api/v1")
	{
		v1.POST("/login", rtr.jwtMiddleware.LoginHandler)
	}
}
