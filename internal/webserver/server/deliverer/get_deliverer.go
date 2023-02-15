package deliverer

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
)

type ResponseOfDeliverer struct {
	types.ModelBase `json:",inline"`

	CreateDelivererReq `json:",inline"`

	Status int `json:"status"`

	Orders []types.OrderForDetail `json:"orders"`
}

func GetDeliverer(ctx *gin.Context) {
	obj := &model.Deliverer{}
	resp := &ResponseOfDeliverer{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func() interface{} { return generateDelivererResponse(obj, resp) },
		LoadAssociationsDBFunc: getDelivererDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getDelivererDetailDB(db *gorm.DB) *gorm.DB {
	return model.GenerateReadDelivererDB(db)
}

func generateDelivererResponse(obj *model.Deliverer, resp *ResponseOfDeliverer) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	for k := range resp.Orders {
		resp.Orders[k].Revise()
	}
	return nil
}
