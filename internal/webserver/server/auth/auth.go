package auth

import (
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	AuthMiddleware.LoginHandler(ctx)
}

func Logout(ctx *gin.Context) {
	AuthMiddleware.LogoutHandler(ctx)
}

func RefreshToken(ctx *gin.Context) {
	AuthMiddleware.RefreshHandler(ctx)
}
