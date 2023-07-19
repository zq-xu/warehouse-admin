package product

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
		TransObjToRespFunc:     func(ac *auth.AccessControl) interface{} { return generateProductResponse(obj, resp, ac) },
		LoadAssociationsDBFunc: getProductDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getProductDetailDB(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	switch auth.RoleSet[ac.User.Role] {
	case auth.UserUserRole:
		return model.GenerateReadProductDB(db, db)
	default:
		return model.GenerateReadProductDB(db, db).Preload("ProductLots.Supplier")
	}
}

func generateProductResponse(obj *model.ProductDetail, resp *ResponseOfProduct, ac *auth.AccessControl) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	if auth.RoleSet[ac.User.Role] == auth.UserUserRole {
		for k := range resp.ProductLots {
			resp.ProductLots[k].SupplierId = ""
			resp.ProductLots[k].Supplier = nil
		}
	}

	return nil
}
