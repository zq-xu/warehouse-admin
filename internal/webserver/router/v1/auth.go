package v1

import (
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/auth"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	authGroupPath = "/auth"
	AuthGroup     = &router.APIGroup{
		RelativePath: authGroupPath,
		APIs: []*router.API{
			// The login method should be ignored by the auth middleware,
			// so not necessary here.
			//{Method: http.MethodPost, Path: "/login", Handler: auth.Login},
			{Method: http.MethodPost, Path: "/refresh_token", Handler: auth.RefreshToken},
			{Method: http.MethodPost, Path: "/logout", Handler: auth.Logout},
		},
	}
)

func init() {
	registerAPIGroup(AuthGroup)
}
