package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func ListProduct(ctx *gin.Context) {
	listObj := make([]model.ProductDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Product{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateProductListResponse(listObj, ac) },
		LoadAssociationsDBFunc: listProductDetailDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func listProductDetailDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadProductDB(db, queryDB).Preload("Category")
}

func generateProductListResponse(objList []model.ProductDetail, ac *auth.AccessControl) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfProduct{}
		_ = generateProductResponse(&v, &r, ac)
		items = append(items, r)
	}

	return items
}
