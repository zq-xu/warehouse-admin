package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/store"
)

var (
	roleSet = map[int]string{
		0: "user",
		1: "admin",
		2: "super",
	}
)

type ResponseOfUserInfo struct {
	types.ModelBase `json:",inline"`

	Name     string `json:"name"`
	UserName string `json:"userName"`
	UserRole string `json:"userRole"`
	Status   int    `json:"status"`
}

func (d *ResponseOfUserInfo) Role(r int) {
	d.UserRole = roleSet[r]
}

func (d *ResponseOfUserInfo) Alias(a string) {
	d.UserName = a
}

func GetUserInfo(ctx *gin.Context) {
	name := ctx.GetString(types.AuthUserNameToken)
	obj, ei := getUserModel(ctx, name)
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

func getUserModel(ctx context.Context, name string) (*model.User, *response.ErrorInfo) {
	db := store.DB(ctx)

	obj := &model.User{}
	err := db.Where("name = ?", name).First(obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCommonError(response.NotFoundErrorCode)
		}

		return nil, response.NewStorageError(response.StorageErrorCode, err)
	}
	return obj, nil
}

func generateUserResponse(obj *model.User) (*ResponseOfUserInfo, *response.ErrorInfo) {
	resp := &ResponseOfUserInfo{}

	err := copier.Copy(resp, obj)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return resp, nil
}
