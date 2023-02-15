package productlot

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/list"
)

type ListResponseOfProductLot []ResponseOfProductLot

func ListProductLot(ctx *gin.Context) {
	listObj := make([]model.ProductLot, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.ProductLot{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func() interface{} { return generateProductLotListResponse(listObj) },
		LoadAssociationsDBFunc: listProductLotDetailDB,
		GenerateQueryFunc:      loadProductLotListQuery,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func listProductLotDetailDB(db, queryDB *gorm.DB) *gorm.DB {
	return queryDB
}

func loadProductLotListQuery(db *gorm.DB, reqParams *list.Params) *gorm.DB {
	productIdQuery := reqParams.Queries[types.ProductIdQuery]
	if productIdQuery != "" {
		db = db.Where("product_id = ?", productIdQuery)
	}
	return db
}

func generateProductLotListResponse(objList []model.ProductLot) interface{} {
	items := make(ListResponseOfProductLot, 0)

	for _, v := range objList {
		r := ResponseOfProductLot{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy product obj to response. %v", err)
			return nil
		}
		items = append(items, r)
	}

	return items
}
