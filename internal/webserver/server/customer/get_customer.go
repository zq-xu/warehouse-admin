package customer

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

type ResponseOfCustomer struct {
	types.ModelBase `json:",inline"`

	CreateCustomerReq `json:",inline"`

	TotalPrice float32 `json:"totalPrice"`
	TotalPaid  float32 `json:"totalPaid"`
	Status     int     `json:"status"`

	Orders []types.OrderForDetail `json:"orders"`
}

func GetCustomer(ctx *gin.Context) {
	obj := &model.CustomerDetail{}
	resp := &ResponseOfCustomer{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func(ac *auth.AccessControl) interface{} { return generateCustomerResponse(obj, resp) },
		LoadAssociationsDBFunc: getCustomerDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getCustomerDetailDB(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadCustomerDB(db, db)
}

func generateCustomerResponse(obj *model.CustomerDetail, resp *ResponseOfCustomer) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	for k := range resp.Orders {
		resp.Orders[k].Revise()
	}
	return nil
}
