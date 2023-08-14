package salesman

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UpdateSalesmanReq struct {
	ID      string
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Comment *string `json:"comment"`
}

func UpdateSalesman(ctx *gin.Context) {
	reqParams := &UpdateSalesmanReq{}
	obj := &model.Salesman{}

	conf := &restapi.UpdateConf{
		UpdateReq:    reqParams,
		ModelObj:     obj,
		OptModelFunc: func(db *gorm.DB) *response.ErrorInfo { return optSalesmanModelForUpdate(reqParams, obj) },
	}

	restapi.Update(ctx, conf)
}

func optSalesmanModelForUpdate(reqParams *UpdateSalesmanReq, obj *model.Salesman) *response.ErrorInfo {
	utils.OptStringPtr(&obj.Name, reqParams.Name)
	utils.OptStringPtr(&obj.Phone, reqParams.Phone)
	utils.OptStringPtr(&obj.Comment, reqParams.Comment)

	return nil
}
