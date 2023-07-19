package supplier

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func ListSupplier(ctx *gin.Context) {
	listObj := make([]model.SupplierDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:              &model.Supplier{},
		ModelObjList:          &listObj,
		FuzzySearchColumnList: []string{"name"},
		AuthControl: restapi.AuthControl{
			AuthValidation: func(ac *auth.AccessControl) bool { return ac.User.Role > 0 },
		},
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateSupplierListResponse(listObj) },
		LoadAssociationsDBFunc: listSupplierDetailDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func listSupplierDetailDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadSupplierDB(db, queryDB)
}

func generateSupplierListResponse(objList []model.SupplierDetail) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfSupplier{}
		_ = generateSupplierResponse(&v, &r)
		items = append(items, r)
	}

	return items
}
