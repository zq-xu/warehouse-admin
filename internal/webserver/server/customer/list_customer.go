package customer

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
)

type ListResponseOfCustomer []ResponseOfCustomer

func ListCustomer(ctx *gin.Context) {
	listObj := make([]model.CustomerDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Customer{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func() interface{} { return generateCustomerListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB) *gorm.DB {
	return model.GenerateReadCustomerDB(queryDB, model.GenerateCustomerAssociationsQuery(db))
}

func generateCustomerListResponse(objList []model.CustomerDetail) interface{} {
	items := make(ListResponseOfCustomer, 0)

	for _, v := range objList {
		r := ResponseOfCustomer{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy customer obj to response. %v", err)
			return nil
		}

		for k := range r.Orders {
			r.Orders[k].Revise()
		}

		items = append(items, r)
	}

	return items
}
