package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func ListCategory(ctx *gin.Context) {
	listObj := make([]model.CategoryDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Category{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		BaseInfoColumnList:     model.CategoryBaseInfoColumns,
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateCategoryListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadCategoryListDB(db, queryDB)
}

func generateCategoryListResponse(objList []model.CategoryDetail) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfCategory{}
		_ = generateCategoryResponse(&v, &r)
		items = append(items, r)
	}

	return items
}
