package salesman

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func ListSalesman(ctx *gin.Context) {
	listObj := make([]model.SalesmanDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Salesman{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		BaseInfoColumnList:     model.SalesmanBaseInfoColumns,
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateSalesmanListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadSalesmanDB(db, queryDB)
}

func generateSalesmanListResponse(objList []model.SalesmanDetail) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfSalesman{}
		_ = generateSalesmanResponse(&v, &r)
		items = append(items, r)
	}

	return items
}
