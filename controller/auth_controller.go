package controller

import (
	"app/configs"
	handler "app/handlers"

	"github.com/gin-gonic/gin"
)

type IAuthController struct {
}
type AuthController struct {
	Gin         *gin.Engine
	AuthHandler handler.IAuthHandler
}

func NewAuthController(a handler.IAuthHandler, ginInsntance *gin.Engine) *AuthController {
	return &AuthController{
		Gin:         ginInsntance,
		AuthHandler: a,
	}
}
func (a *AuthController) AuthEndpoints() {
	a.Gin.GET(configs.API("auth"), a.AuthHandler.GoogleSSOHandler)
}
