package product

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
)

type ResponseOfProduct struct {
	types.ModelBase `json:",inline"`

	CreateProductReq `json:",inline"`

	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`

	TotalCount int `json:"totalCount"`
	SoldCount  int `json:"soldCount"`
	Stocks     int `json:"stocks"`

	Status int `json:"status"`

	ProductLots []types.ProductLotForDetail `json:"productLots"`
}

func GetProduct(ctx *gin.Context) {
	obj := &model.ProductDetail{}
	resp := &ResponseOfProduct{}

	conf := &restapi.DetailConf{
		ModelObj:               obj,
		RespObj:                resp,
		TransObjToRespFunc:     func() interface{} { return generateProductResponse(obj, resp) },
		LoadAssociationsDBFunc: getProductDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getProductDetailDB(db *gorm.DB) *gorm.DB {
	return model.GenerateReadProductDB(db, db).
		Preload("ProductLots.Supplier")
}

func generateProductResponse(obj *model.ProductDetail, resp *ResponseOfProduct) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return nil
}
