package deliverer

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func ListDeliverer(ctx *gin.Context) {
	listObj := make([]model.Deliverer, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Deliverer{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		BaseInfoColumnList:     model.DelivererBaseInfoColumns,
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateDelivererListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadDelivererDB(queryDB)
}

func generateDelivererListResponse(objList []model.Deliverer) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfDeliverer{}

		_ = generateDelivererResponse(&v, &r)

		items = append(items, r)
	}

	return items
}
