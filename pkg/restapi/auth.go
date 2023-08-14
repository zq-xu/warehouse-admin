package restapi

import (
	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

type AuthControl struct {
	AccessControl  *auth.AccessControl
	AuthValidation func(ac *auth.AccessControl) bool
}

func (ac *AuthControl) Validate(ctx *gin.Context) *response.ErrorInfo {
	a, ei := auth.GetAccessControl(ctx, ctx.GetString(auth.AuthUserIDToken))
	if ei != nil {
		return ei
	}

	ac.AccessControl = a

	if ac.AuthValidation != nil && !ac.AuthValidation(a) {
		return response.NewCommonError(response.InvalidAuthErrorCode)
	}

	return nil
}
