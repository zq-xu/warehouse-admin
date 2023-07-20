package category

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

type ResponseOfCategory struct {
	types.CategoryForDetail `json:",inline"`

	Products []types.ProductForDetail `json:"products"`

	ProductCount int `json:"productCount"`
}

func GetCategory(ctx *gin.Context) {
	obj := &model.CategoryDetail{}
	resp := &ResponseOfCategory{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func(ac *auth.AccessControl) interface{} { return generateCategoryResponse(obj, resp) },
		LoadAssociationsDBFunc: getCategoryDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getCategoryDetailDB(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadCategoryDetailDB(db)
}

func generateCategoryResponse(obj *model.CategoryDetail, resp *ResponseOfCategory) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return nil
}
