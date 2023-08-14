package auth

import (
	"context"

	"zq-xu/warehouse-admin/pkg/restapi/response"
)

type AccessControl struct {
	User *User
}

func GetAccessControl(ctx context.Context, id string) (*AccessControl, *response.ErrorInfo) {
	u, ei := GetUserModel(ctx, id)
	if ei != nil {
		return nil, ei
	}

	return &AccessControl{User: u}, nil
}
