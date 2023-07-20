package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UpdateCategoryReq struct {
	ID   string
	Name *string `json:"name"`
}

func UpdateCategory(ctx *gin.Context) {
	reqParams := &UpdateCategoryReq{}
	obj := &model.Category{}

	conf := &restapi.UpdateConf{
		UpdateReq:    reqParams,
		ModelObj:     obj,
		OptModelFunc: func(db *gorm.DB) *response.ErrorInfo { return optCategoryModelForUpdate(reqParams, obj) },
	}

	restapi.Update(ctx, conf)
}

func optCategoryModelForUpdate(reqParams *UpdateCategoryReq, obj *model.Category) *response.ErrorInfo {
	utils.OptStringPtr(&obj.Name, reqParams.Name)
	return nil
}
