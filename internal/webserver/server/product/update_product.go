package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"zq-xu/warehouse-admin/pkg/store"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UpdateProductReq struct {
	ID             string
	Name           *string  `json:"name"`
	Price          *float32 `json:"price"`
	StorageAddress *string  `json:"storageAddress"`
	Comment        *string  `json:"comment"`
	CategoryId     *string  `json:"categoryId"`
}

func UpdateProduct(ctx *gin.Context) {
	reqParams := &UpdateProductReq{}
	obj := &model.Product{}

	conf := &restapi.UpdateConf{
		UpdateReq:    reqParams,
		ModelObj:     obj,
		OptModelFunc: func(db *gorm.DB) *response.ErrorInfo { return optProductModelForUpdate(db, reqParams, obj) },
	}

	restapi.Update(ctx, conf)
}

func optProductModelForUpdate(db *gorm.DB, reqParams *UpdateProductReq, obj *model.Product) *response.ErrorInfo {
	ei := optCategoryID(db, reqParams, obj)
	if ei != nil {
		return ei
	}

	utils.OptStringPtr(&obj.Name, reqParams.Name)
	utils.OptFloat32Ptr(&obj.Price, reqParams.Price)
	utils.OptStringPtr(&obj.StorageAddress, reqParams.StorageAddress)
	utils.OptStringPtr(&obj.Comment, reqParams.Comment)

	return nil
}

func optCategoryID(db *gorm.DB, reqParams *UpdateProductReq, obj *model.Product) *response.ErrorInfo {
	if reqParams.CategoryId == nil {
		return nil
	}

	ei := store.EnsureExistByID(db, &model.Category{}, *reqParams.CategoryId)
	if ei != nil {
		return ei
	}

	obj.CategoryID = utils.GetInt64PtrByStringPtrDefaultNil(reqParams.CategoryId)
	return nil
}
