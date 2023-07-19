package supplier

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UpdateSupplierReq struct {
	ID          string
	Name        *string `json:"name"`
	Phone       *string `json:"phone"`
	Address     *string `json:"address"`
	Comment     *string `json:"comment"`
	InvoiceInfo *string `json:"invoiceInfo"`
}

func UpdateSupplier(ctx *gin.Context) {
	reqParams := &UpdateSupplierReq{}
	obj := &model.Supplier{}

	conf := &restapi.UpdateConf{
		UpdateReq: reqParams,
		ModelObj:  obj,
		AuthControl: restapi.AuthControl{
			AuthValidation: func(ac *auth.AccessControl) bool { return ac.User.Role > 0 },
		},
		OptModelFunc: func(db *gorm.DB) *response.ErrorInfo { return optSupplierModelForUpdate(reqParams, obj) },
	}

	restapi.Update(ctx, conf)
}

func optSupplierModelForUpdate(reqParams *UpdateSupplierReq, obj *model.Supplier) *response.ErrorInfo {
	utils.OptStringPtr(&obj.Name, reqParams.Name)
	utils.OptStringPtr(&obj.Phone, reqParams.Phone)
	utils.OptStringPtr(&obj.Address, reqParams.Address)
	utils.OptStringPtr(&obj.Comment, reqParams.Comment)
	utils.OptStringPtr(&obj.InvoiceInfo, reqParams.InvoiceInfo)

	return nil
}
