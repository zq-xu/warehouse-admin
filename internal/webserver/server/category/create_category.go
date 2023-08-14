package category

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/store"
)

type CreateCategoryReq struct {
	Name string `json:"name"`
}

func CreateCategory(ctx *gin.Context) {
	reqParams, ei := newCreateCategoryReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	_, ei = auth.GetAccessControl(ctx, ctx.GetString(auth.AuthUserIDToken))
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		ei = store.EnsureNotExistByName(db, &model.Category{}, reqParams.Name)
		if ei != nil {
			return ei
		}

		obj, ei := generateCategoryModelForCreation(reqParams)
		if ei != nil {
			return ei
		}

		err := store.Create(db, obj)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		log.Logger.Infof("Succeed to create category %d/%s", obj.ID, obj.Name)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func newCreateCategoryReq(ctx *gin.Context) (*CreateCategoryReq, *response.ErrorInfo) {
	reqBody := &CreateCategoryReq{}
	err := ctx.ShouldBindJSON(reqBody)
	if err != nil {
		return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
	}

	return reqBody, nil
}

func generateCategoryModelForCreation(reqParams *CreateCategoryReq) (*model.Category, *response.ErrorInfo) {
	t := &model.Category{
		Model: store.GenerateModel(),
	}

	err := copier.Copy(t, reqParams)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return t, nil
}
