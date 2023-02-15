package product

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
)

type ListResponseOfProduct []ResponseOfProduct

func ListProduct(ctx *gin.Context) {
	listObj := make([]model.ProductDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Product{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func() interface{} { return generateProductListResponse(listObj) },
		LoadAssociationsDBFunc: listProductDetailDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func listProductDetailDB(db, queryDB *gorm.DB) *gorm.DB {
	return model.GenerateReadProductDB(queryDB, model.GenerateProductAssociationsQuery(db))
}

func generateProductListResponse(objList []model.ProductDetail) interface{} {
	items := make(ListResponseOfProduct, 0)

	for _, v := range objList {
		r := ResponseOfProduct{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy product obj to response. %v", err)
			return nil
		}
		items = append(items, r)
	}

	return items
}
