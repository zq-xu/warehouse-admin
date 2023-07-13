package order

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

type ListResponseOfOrder []types.OrderForDetail

func ListOrder(ctx *gin.Context) {
	listObj := make([]model.OrderDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Order{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func() interface{} { return generateOrderListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
		GenerateQueryFunc:      loadOrderListQuery,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB) *gorm.DB {
	return model.GenerateReadOrderDB(db, queryDB)
}

func loadOrderListQuery(db *gorm.DB, reqParams *list.Params) *gorm.DB {
	customerId := reqParams.Queries[types.CustomerIdQuery]
	if customerId != "" {
		db = db.Where("customer_id = ?", customerId)
	}

	salesmanId := reqParams.Queries[types.SalesmanIdQuery]
	if salesmanId != "" {
		db = db.Where("salesman_id = ?", salesmanId)
	}

	delivererId := reqParams.Queries[types.DelivererIdQuery]
	if delivererId != "" {
		db = db.Where("deliverer_id = ?", delivererId)
	}

	return db
}

func generateOrderListResponse(objList []model.OrderDetail) interface{} {
	items := make(ListResponseOfOrder, 0)

	for _, v := range objList {
		r := types.OrderForDetail{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy order obj to response. %v", err)
			return nil
		}

		items = append(items, r)
	}

	return items
}
