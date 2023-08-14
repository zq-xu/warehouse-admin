package order

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

func GetOrder(ctx *gin.Context) {
	obj := &model.OrderDetail{}
	resp := &types.OrderForDetail{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func(ac *auth.AccessControl) interface{} { return generateOrderResponse(obj, resp) },
		LoadAssociationsDBFunc: getOrderDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getOrderDetailDB(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadOrderDB(db, db)
}

func generateOrderResponse(obj *model.OrderDetail, resp *types.OrderForDetail) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return nil
}
