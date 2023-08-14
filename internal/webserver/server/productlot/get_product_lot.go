package productlot

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/brick"
	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func GetProductLot(ctx *gin.Context) {
	obj := &model.ProductLot{}
	resp := &types.ProductLotForDetail{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func(ac *auth.AccessControl) interface{} { return generateProductLotResponse(obj, resp, ac) },
		LoadAssociationsDBFunc: getProductLotDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getProductLotDetailDB(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return brick.OptProductLotDBByAuth(db, ac)
}

func generateProductLotResponse(obj *model.ProductLot, resp *types.ProductLotForDetail, ac *auth.AccessControl) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	brick.OptProductLotRespByAuth(resp, ac)
	return nil
}
