package product

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi"
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

func (erd *ExcelRowData) ID(id interface{}) {
	erd.Id = fmt.Sprintf("%v", id)
}

func (erd *ExcelRowData) ProductLots(list []model.ProductLot) {
	bs, _ := json.Marshal(&list)
	erd.ProductLotsJson = string(bs)
}

func ExportProduct(ctx *gin.Context) {
	listObj := make([]model.Product, 0)
	data := make([]interface{}, 0)

	conf := &restapi.ListConf{
		ModelObj:              &model.Product{},
		ModelObjList:          &listObj,
		FuzzySearchColumnList: []string{"url"},
		TransObjToRespFunc: func() interface{} {
			data = generateProductListExcelData(listObj, data)
			return data
		},
		LoadAssociationsDBFunc: exportProductListDB,
		ResponseWriteFunc:      func(ctx *gin.Context) { writeProductToExcel(ctx, data) },
	}

	restapi.ApiListInstance.List(ctx, conf)
}

func exportProductListDB(db, queryDB *gorm.DB) *gorm.DB {
	return model.GenerateReadProductDB(queryDB, db).
		Preload("ProductLots")
}

func generateProductListExcelData(objList []model.Product, data []interface{}) []interface{} {
	for _, v := range objList {
		r := ExcelRowData{}

		err := copier.Copy(&r, &v)
		if err != nil {
			log.Logger.Errorf("Failed to copy product obj to response. %v", err)
			return nil
		}

		data = append(data, r)
	}

	return data
}

func writeProductToExcel(ctx *gin.Context, data []interface{}) {
	ex := utils.NewExcelExport("")
	ex.WriteExcelByStruct("", excelTittleList, data)
	_ = ex.ExportExcelToGin(productExcelName, ctx)
}
