package salesman

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
)

type ListResponseOfSalesman []ResponseOfSalesman

func ListSalesman(ctx *gin.Context) {
	listObj := make([]model.SalesmanDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Salesman{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func() interface{} { return generateSalesmanListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB) *gorm.DB {
	return model.GenerateReadSalesmanDB(queryDB, model.GenerateSalesmanAssociationsQuery(db))
}

func generateSalesmanListResponse(objList []model.SalesmanDetail) interface{} {
	items := make(ListResponseOfSalesman, 0)

	for _, v := range objList {
		r := ResponseOfSalesman{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy salesman obj to response. %v", err)
			return nil
		}

		for k := range r.Orders {
			r.Orders[k].Revise()
		}

		items = append(items, r)
	}

	return items
}