package order

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
)

type ResponseOfOrder struct {
	types.ModelBase `json:",inline"`

	CreateOrderReq `json:",inline"`

	OrderNo    string  `json:"orderNo"`
	TotalPrice float32 `json:"totalPrice"`
	TotalPaid  float32 `json:"totalPaid"`
	Status     int     `json:"status"`

	Customer      types.CustomerForDetail       `json:"customer"`
	OrderProducts []types.OrderProductForDetail `json:"products"`
	Salesman      types.SalesmanForDetail       `json:"salesman"`
	Deliverer     types.DelivererForDetail      `json:"deliverer"`
}

func GetOrder(ctx *gin.Context) {
	obj := &model.OrderDetail{}
	resp := &ResponseOfOrder{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func() interface{} { return generateSupplierResponse(obj, resp) },
		LoadAssociationsDBFunc: getSupplierDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getSupplierDetailDB(db *gorm.DB) *gorm.DB {
	return model.GenerateReadOrderDB(db, model.GenerateOrderAssociationsQuery(db))
}

func generateSupplierResponse(obj *model.OrderDetail, resp *ResponseOfOrder) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return nil
}
