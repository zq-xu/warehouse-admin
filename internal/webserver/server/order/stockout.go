package order

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/store"
)

type StockOutReq struct {
	OrderID     string               `json:"orderId"`
	ProductLots []StockOutProductLot `json:"productLots"`

	Order model.Order `json:"-" copier:"-"`
}

type StockOutProductLot struct {
	ProductLotID string `json:"productLotId"`
	Count        int    `json:"count"`
	Comment      string `json:"comment"`

	ProductLot model.ProductLot `json:"-" copier:"-"`
}

// StockOut
/*
Key points:
 1. the productLots are bound to the order product or not;
 2. the stockOut product can be updated with the stockOutID.
*/
func StockOut(ctx *gin.Context) {
	reqParams, ei := newStockOutProductReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	_, ei = auth.GetAccessControl(ctx, ctx.GetString(auth.AuthUserIDToken))
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		err := store.GetByID(db, &reqParams.Order, reqParams.OrderID)
		if err != nil {
			return response.NewCommonError(response.NotFoundErrorCode)
		}

		for _, v := range reqParams.ProductLots {
			err = store.GetByID(db, &v.ProductLot, v.ProductLotID)
			if ei != nil {
				return response.NewCommonError(response.NotFoundErrorCode)
			}

			obj := generateStockOutModel(reqParams.Order.ID, &v)

			err = store.Create(db, obj)
			if err != nil {
				return response.NewStorageError(response.StorageErrorCode, err)
			}
		}
		log.Logger.Infof("Succeed to stock out order %d", reqParams.Order.ID)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func newStockOutProductReq(ctx *gin.Context) (*StockOutReq, *response.ErrorInfo) {
	reqBody := &StockOutReq{}
	err := ctx.ShouldBindJSON(reqBody)
	if err != nil {
		return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
	}

	return reqBody, nil
}

func generateStockOutModel(orderId int64, pl *StockOutProductLot) *model.StockOut {
	return &model.StockOut{
		Model:        store.GenerateModel(),
		OrderID:      orderId,
		ProductLotID: pl.ProductLot.ID,
		Count:        pl.Count,
		Comment:      pl.Comment,
	}
}
