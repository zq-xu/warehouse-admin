package supplier

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
)

type ListResponseOfSupplier []ResponseOfSupplier

func ListSupplier(ctx *gin.Context) {
	listObj := make([]model.SupplierDetail, 0)

	conf := &restapi.ListConf{
		ModelObj:               &model.Supplier{},
		ModelObjList:           &listObj,
		FuzzySearchColumnList:  []string{"name"},
		TransObjToRespFunc:     func() interface{} { return generateSupplierListResponse(listObj) },
		LoadAssociationsDBFunc: listSupplierDetailDB,
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func listSupplierDetailDB(db, queryDB *gorm.DB) *gorm.DB {
	return model.GenerateReadSupplierDB(queryDB, model.GenerateSupplierAssociationsQuery(db))
}

func generateSupplierListResponse(objList []model.SupplierDetail) interface{} {
	items := make(ListResponseOfSupplier, 0)

	for _, v := range objList {
		r := ResponseOfSupplier{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy supplier obj to response. %v", err)
			return nil
		}
		items = append(items, r)
	}

	return items
}
