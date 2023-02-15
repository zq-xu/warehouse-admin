package salesman

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
)

type GetSalesmanReq struct {
	ID string
}

type ResponseOfSalesman struct {
	types.ModelBase `json:",inline"`

	CreateSalesmanReq `json:",inline"`

	TotalPrice float32 `json:"totalPrice"`
	TotalPaid  float32 `json:"totalPaid"`
	Status     int     `json:"status"`

	Orders []types.OrderForDetail `json:"orders"`
}

func GetSalesman(ctx *gin.Context) {
	obj := &model.SalesmanDetail{}
	resp := &ResponseOfSalesman{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func() interface{} { return generateSupplierResponse(obj, resp) },
		LoadAssociationsDBFunc: getSupplierDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getSupplierDetailDB(db *gorm.DB) *gorm.DB {
	return model.GenerateReadSalesmanDB(db, model.GenerateSalesmanAssociationsQuery(db))
}

func generateSupplierResponse(obj *model.SalesmanDetail, resp *ResponseOfSalesman) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	for k := range resp.Orders {
		resp.Orders[k].Revise()
	}
	return nil
}
