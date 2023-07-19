package product

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/utils"
)

const productExcelName = "product"

var excelTittleList []string

type ExcelRowData struct {
	Id        string
	CreatedAt utils.UnixTime
	UpdatedAt utils.UnixTime

	Name           string
	Price          float32
	StorageAddress string
	Comment        string

	Image     string
	Thumbnail string

	TotalCount int
	Status     int

	ProductLotsJson string
}

func init() {
	excelTittleList = make([]string, 0)

	t := reflect.TypeOf(ExcelRowData{})
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		excelTittleList = append(excelTittleList, t.Field(i).Name)
	}
}

func (erd *ExcelRowData) ProductLots(list []types.ProductLotForDetail) {
	bs, _ := json.Marshal(&list)
	erd.ProductLotsJson = string(bs)
}

func ExportProduct(ctx *gin.Context) {
	listObj := make([]model.ProductDetail, 0)
	data := make([]interface{}, 0)

	conf := &restapi.ListConf{
		ModelObj:              &model.Product{},
		ModelObjList:          &listObj,
		FuzzySearchColumnList: []string{"url"},
		TransObjToRespFunc: func(ac *auth.AccessControl) []interface{} {
			data = generateProductListExcelData(listObj, data, ac)
			return data
		},
		LoadAssociationsDBFunc: exportProductListDB,
		ResponseWriteFunc:      func(ctx *gin.Context) { writeProductToExcel(ctx, data) },
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func exportProductListDB(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB {
	return model.GenerateReadProductDB(db, queryDB).Preload("ProductLots")
}

func generateProductListExcelData(objList []model.ProductDetail, data []interface{}, ac *auth.AccessControl) []interface{} {
	for _, v := range objList {
		r := ResponseOfProduct{}
		_ = generateProductResponse(&v, &r, ac)
		erd := &ExcelRowData{}
		_ = copier.Copy(erd, &r)
		data = append(data, *erd)
	}
	return data
}

func writeProductToExcel(ctx *gin.Context, data []interface{}) {
	ex := utils.NewExcelExport("")
	ex.WriteExcelByStruct("", excelTittleList, data)
	_ = ex.ExportExcelToGin(productExcelName, ctx)
}
