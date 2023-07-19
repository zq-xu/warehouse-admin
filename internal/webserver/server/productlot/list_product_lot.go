package productlot

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/list"
	"zq-xu/warehouse-admin/pkg/router/auth"
)

func ListProductLot(ctx *gin.Context) {
	listObj := make([]model.ProductLot, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.ProductLot{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func(ac *auth.AccessControl) []interface{} { return generateProductLotListResponse(listObj) },
		LoadAssociationsDBFunc: listProductLotDetailDB,
		GenerateQueryFunc:      loadProductLotListQuery,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func listProductLotDetailDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	switch auth.RoleSet[ac.User.Role] {
	case auth.UserUserRole:
		return queryDB.Omit("SupplierID", "Supplier")
	default:
		return queryDB
	}
}

func loadProductLotListQuery(db *gorm.DB, reqParams *list.Params) *gorm.DB {
	productIdQuery := reqParams.Queries[types.ProductIdQuery]
	if productIdQuery != "" {
		db = db.Where("product_id = ?", productIdQuery)
	}
	return db
}

func generateProductLotListResponse(objList []model.ProductLot) []interface{} {
	items := make([]interface{}, 0)

	for _, v := range objList {
		r := ResponseOfProductLot{}
		_ = generateProductLotResponse(&v, &r)
		items = append(items, r)
	}

	return items
}
