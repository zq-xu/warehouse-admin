package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/store"
)

const (
	UserUserRole  UserRole = "user"
	AdminUserRole UserRole = "admin"
	SuperUserRole UserRole = "super"
)

var (
	RoleSet = map[int]UserRole{
		0: UserUserRole,
		1: AdminUserRole,
		2: SuperUserRole,
	}
)

type UserRole string

type ResponseOfUserInfo struct {
	Name     string   `json:"name"`
	UserName string   `json:"userName"`
	UserRole UserRole `json:"userRole"`
	Status   int      `json:"status"`
}

func (r *ResponseOfUserInfo) Role(rl int) {
	r.UserRole = RoleSet[rl]
}

func (r *ResponseOfUserInfo) Alias(a string) {
	r.UserName = a
}

func GetUserInfo(ctx *gin.Context) {
	id := ctx.GetString(AuthUserIDToken)

	obj, ei := GetUserModel(ctx, id)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	res, ei := generateUserResponse(obj)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func GetUserModel(ctx context.Context, id string) (*User, *response.ErrorInfo) {
	db := store.DB(ctx)

	obj := &User{}
	err := db.Where("id = ?", id).First(obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCommonError(response.InvalidAuthErrorCode)
		}

		return nil, response.NewStorageError(response.StorageErrorCode, err)
	}
	return obj, nil
}

func generateUserResponse(obj *User) (*ResponseOfUserInfo, *response.ErrorInfo) {
	resp := &ResponseOfUserInfo{}

	err := copier.Copy(resp, obj)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return resp, nil
}
