package productlot

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/store"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UpdateProductLotReq struct {
	ID              string
	SupplierID      *string         `json:"supplierId"`
	PurchasingPrice *float32        `json:"purchasingPrice"`
	PurchasingDate  *utils.UnixTime `json:"purchasingDate"`
	Paid            *float32        `json:"paid"`
	Count           *int            `json:"count"`
	StorageAddress  *string         `json:"storageAddress"`
	Comment         *string         `json:"comment"`
}

func UpdateProductLot(ctx *gin.Context) {
	reqParams := &UpdateProductLotReq{}
	obj := &model.ProductLot{}

	conf := &restapi.UpdateConf{
		UpdateReq:    reqParams,
		ModelObj:     obj,
		OptModelFunc: func(db *gorm.DB) *response.ErrorInfo { return optProductLotModelForUpdate(db, reqParams, obj) },
	}

	restapi.Update(ctx, conf)
}

func optProductLotModelForUpdate(db *gorm.DB, reqParams *UpdateProductLotReq, obj *model.ProductLot) *response.ErrorInfo {
	ei := optSupplierID(db, reqParams, obj)
	if ei != nil {
		return ei
	}

	utils.OptInt64ByStringPtr(&obj.SupplierID, reqParams.SupplierID)
	utils.OptFloat32Ptr(&obj.PurchasingPrice, reqParams.PurchasingPrice)
	utils.OptTimeByUnixTimePtr(&obj.PurchasingDate, reqParams.PurchasingDate)
	utils.OptFloat32Ptr(&obj.Paid, reqParams.Paid)
	utils.OptIntPtr(&obj.Count, reqParams.Count)
	utils.OptStringPtr(&obj.StorageAddress, reqParams.StorageAddress)
	utils.OptStringPtr(&obj.Comment, reqParams.Comment)

	return nil
}

func optSupplierID(db *gorm.DB, reqParams *UpdateProductLotReq, obj *model.ProductLot) *response.ErrorInfo {
	if reqParams.SupplierID == nil {
		return nil
	}

	ei := store.EnsureExistByID(db, &model.Supplier{}, *reqParams.SupplierID)
	if ei != nil {
		return ei
	}

	utils.OptInt64ByStringPtr(&obj.SupplierID, reqParams.SupplierID)
	return nil
}
