package router

import (
	"zq-xu/warehouse-admin/internal/webserver/router/v1"
	"zq-xu/warehouse-admin/pkg/router"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

const (
	// MaxMultipartMemory 100M
	MaxMultipartMemory = 100 << 20
)

func StartRouter() error {
	r := router.DefaultRouter()

	r.MaxMultipartMemory = MaxMultipartMemory

	// The login method here is to be ignored by the auth middleware
	r.POST("/login", auth.Login)
	r.POST("/refresh_token", auth.RefreshToken)
	r.POST("/logout", auth.Logout)

	v1.Register()

	go router.StartPprof()
	return router.StartRouter(r)
}
