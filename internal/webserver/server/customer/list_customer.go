package customer

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func ListCustomer(ctx *gin.Context) {
	listObj := make([]model.CustomerDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Customer{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		BaseInfoColumnList:     model.CustomerBaseInfoColumns,
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateCustomerListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadCustomerDB(db, queryDB)
}

func generateCustomerListResponse(objList []model.CustomerDetail) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfCustomer{}
		_ = generateCustomerResponse(&v, &r)
		items = append(items, r)
	}

	return items
}
