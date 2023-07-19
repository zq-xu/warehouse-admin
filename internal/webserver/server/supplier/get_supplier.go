package supplier

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

type ResponseOfSupplier struct {
	types.ModelBase `json:",inline"`

	CreateSupplierReq `json:",inline"`

	TotalPrice float32 `json:"totalPrice"`
	TotalPaid  float32 `json:"totalPaid"`

	Status int `json:"status"`

	ProductLots []types.ProductLotForDetail `json:"productLots"`
}

func GetSupplier(ctx *gin.Context) {
	obj := &model.SupplierDetail{}
	resp := &ResponseOfSupplier{}

	conf := &restapi.DetailConf{
		ModelObj: obj,
		RespObj:  resp,
		AuthControl: restapi.AuthControl{
			AuthValidation: func(ac *auth.AccessControl) bool { return ac.User.Role > 0 },
		},
		TransObjToRespFunc:     func(ac *auth.AccessControl) interface{} { return generateSupplierResponse(obj, resp) },
		LoadAssociationsDBFunc: getSupplierDetailDB,
	}

	restapi.GetDetail(ctx, conf)
}

func getSupplierDetailDB(db *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadSupplierDB(db, db)
}

func generateSupplierResponse(obj *model.SupplierDetail, resp *ResponseOfSupplier) *response.ErrorInfo {
	err := copier.Copy(resp, obj)
	if err != nil {
		return response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return nil
}
