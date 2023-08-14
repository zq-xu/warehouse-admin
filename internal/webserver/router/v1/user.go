package v1

import (
	"fmt"
	"net/http"
	"zq-xu/warehouse-admin/pkg/router/auth"

	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	userGroupPath = "/users"
	userPath      = fmt.Sprintf("/:%s", types.IDParam)
	UserGroup     = &router.APIGroup{
		RelativePath: userGroupPath,
		APIs: []*router.API{
			{Method: http.MethodGet, Path: "/info", Handler: auth.GetUserInfo},
		},
	}
)

func init() {
	registerAPIGroup(UserGroup)
}
