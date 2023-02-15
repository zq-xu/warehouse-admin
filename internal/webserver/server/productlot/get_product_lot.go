package productlot

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
)

type ResponseOfProductLot struct {
	types.ProductLotForDetail `json:",inline"`
}

func GetProductLot(ctx *gin.Context) {
	obj := &model.ProductLot{}
	resp := &ResponseOfProductLot{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func() interface{} { return generateProductLotResponse(obj, resp) },
		LoadAssociationsDBFunc: getProductLotDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getProductLotDetailDB(db *gorm.DB) *gorm.DB {
	return db
}

func generateProductLotResponse(obj *model.ProductLot, resp *ResponseOfProductLot) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return nil
}
