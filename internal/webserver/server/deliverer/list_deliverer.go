package deliverer

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
)

type ListResponseOfDeliverer []ResponseOfDeliverer

func ListDeliverer(ctx *gin.Context) {
	listObj := make([]model.Deliverer, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Deliverer{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func() interface{} { return generateDelivererListResponse(listObj) },
		LoadAssociationsDBFunc: LoadAssociationsDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func LoadAssociationsDB(db, queryDB *gorm.DB) *gorm.DB {
	return model.GenerateReadDelivererDB(queryDB)
}

func generateDelivererListResponse(objList []model.Deliverer) interface{} {
	items := make(ListResponseOfDeliverer, 0)

	for _, v := range objList {
		r := ResponseOfDeliverer{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy deliverer obj to response. %v", err)
			return nil
		}

		for k := range r.Orders {
			r.Orders[k].Revise()
		}

		items = append(items, r)
	}

	return items
}
