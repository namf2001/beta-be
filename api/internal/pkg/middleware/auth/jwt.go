package auth

import (
	"net/http"
	"time"

	"beta-be/internal/model"
	"beta-be/internal/pkg/config"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	UserNameKey = "username"
	EmailKey    = "email"
)

func NewJWTAuth() (*jwt.GinJWTMiddleware, error) {
	timeout := time.Hour
	if config.AppConfig.Environment == config.EnvironmentDevelopment {
		timeout = time.Duration(876010) * time.Hour
	} else {
		if config.AppConfig.JWTExpired != 0 {
			timeout = time.Duration(config.AppConfig.JWTExpired) * time.Second
		}
	}
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.AppConfig.JWTRealm,
		Key:             []byte(config.AppConfig.JWTSecret),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		PayloadFunc:     PayloadFunc,
		IdentityHandler: IdentityHandler,
		Authorizator:    Authorize,
		Unauthorized:    Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(model.User)
		return jwt.MapClaims{
			jwt.IdentityKey: u.ID,
			UserNameKey:     u.UserName,
			EmailKey:        u.Email,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey": claims["identity"],
		"UserName":    claims["username"],
		"Email":       claims["email"],
		"UserId":      claims["identity"],
	}
}

func Authorize(data interface{}, c *gin.Context) bool {

	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(model.User)

		c.Set("userId", u.ID)
		c.Set("userName", u.UserName)
		c.Set("email", u.Email)
		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
