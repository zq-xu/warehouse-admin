package customer

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UpdateCustomerReq struct {
	ID          string
	Name        *string `json:"name"`
	Phone       *string `json:"phone"`
	Address     *string `json:"address"`
	Comment     *string `json:"comment"`
	InvoiceInfo *string `json:"invoiceInfo"`
}

func UpdateCustomer(ctx *gin.Context) {
	reqParams := &UpdateCustomerReq{}
	obj := &model.Customer{}

	conf := &restapi.UpdateConf{
		UpdateReq:    reqParams,
		ModelObj:     obj,
		OptModelFunc: func(db *gorm.DB) *response.ErrorInfo { return optCustomerModelForUpdate(reqParams, obj) },
	}

	restapi.Update(ctx, conf)
}

func optCustomerModelForUpdate(reqParams *UpdateCustomerReq, obj *model.Customer) *response.ErrorInfo {
	utils.OptStringPtr(&obj.Name, reqParams.Name)
	utils.OptStringPtr(&obj.Phone, reqParams.Phone)
	utils.OptStringPtr(&obj.Address, reqParams.Address)
	utils.OptStringPtr(&obj.Comment, reqParams.Comment)
	utils.OptStringPtr(&obj.InvoiceInfo, reqParams.InvoiceInfo)

	return nil
}
