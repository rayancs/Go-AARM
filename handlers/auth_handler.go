package handler

import (
	"app/configs"
	"app/services"
	"app/types"

	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	// handler
	GoogleSSOHandler(g *gin.Context)
}
type AuthHandler struct {
	AuthServ services.IAuthService
}

func NewAuthHandler(x services.IAuthService) *AuthHandler {
	return &AuthHandler{
		AuthServ: x,
	}
}

func (a *AuthHandler) GoogleSSOHandler(g *gin.Context) {
	//
	code := g.Query("code")
	// validation for code
	if code == "" {
		errJson := types.NewHttpResponse(configs.Oops, "code_missing", 400, nil)
		g.JSON(errJson.Status, errJson)
		return
	}
	// pass it to auth funtion
	res := a.AuthServ.GoogleSSO(code)
	g.JSON(res.Status, res)
}
