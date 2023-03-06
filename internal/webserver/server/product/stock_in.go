package product

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/store"
)

type StockInProductReq struct {
	ProductID       string    `json:"productId"`
	SupplierID      string    `json:"supplierId"`
	Count           int       `json:"count"`
	PurchasingPrice float32   `json:"purchasingPrice"`
	Paid            float32   `json:"paid"`
	PurchasingDate  time.Time `json:"purchasingDate"`
	LotNo           string    `json:"lotNo"`
	StorageAddress  string    `json:"storageAddress"`
	Comment         string    `json:"comment"`

	Product  model.Product  `json:"-" copier:"-"`
	Supplier model.Supplier `json:"-" copier:"-"`
}

func StockInProduct(ctx *gin.Context) {
	reqParams, ei := newStockInProductReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		err := store.GetByID(db, &reqParams.Product, reqParams.ProductID)
		if err != nil {
			return response.NewCommonError(response.NotFoundErrorCode)
		}

		err = store.GetByID(db, &reqParams.Supplier, reqParams.SupplierID)
		if ei != nil {
			return response.NewCommonError(response.NotFoundErrorCode)
		}

		obj, ei := generateProductLotModel(reqParams)
		if ei != nil {
			return ei
		}

		err = store.Create(db, obj)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		log.Logger.Infof("Succeed to stock in product %d/%s", obj.ID, obj.ProductID)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func newStockInProductReq(ctx *gin.Context) (*StockInProductReq, *response.ErrorInfo) {
	reqBody := &StockInProductReq{}
	err := ctx.ShouldBindJSON(reqBody)
	if err != nil {
		return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
	}

	return reqBody, nil
}

func generateProductLotModel(reqParams *StockInProductReq) (*model.ProductLot, *response.ErrorInfo) {
	t := &model.ProductLot{
		Model: store.GenerateModel(),
	}

	err := copier.Copy(t, reqParams)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	t.ProductID = reqParams.Product.ID
	t.SupplierID = reqParams.Supplier.ID

	return t, nil
}
