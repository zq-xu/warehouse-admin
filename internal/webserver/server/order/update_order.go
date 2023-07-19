package order

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/restapi"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/store"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UpdateOrderReq struct {
	CustomerId      *string                      `json:"customerId"`
	SalesmanId      *string                      `json:"salesmanId"`
	DelivererId     *string                      `json:"delivererId"`
	PayMode         *int                         `json:"payMode"`
	Paid            *float32                     `json:"paid"`
	Comment         *string                      `json:"comment"`
	DeliveryMode    *int                         `json:"deliveryMode"`
	DeliveryAt      *utils.UnixTime              `json:"deliveryAt"`
	DeliveryAddress *string                      `json:"deliveryAddress"`
	Products        map[string]*ProductForUpdate `json:"products"`
}

type ProductForUpdate struct {
	Count      *int     `json:"count"`
	Paid       *float32 `json:"paid"`
	FinalPrice *float32 `json:"finalPrice"`
	Comment    *string  `json:"comment"`
}

func UpdateOrder(ctx *gin.Context) {
	_, ei := auth.GetAccessControl(ctx, ctx.GetString(auth.AuthUserIDToken))
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	reqParams := &UpdateOrderReq{}
	obj := &model.Order{}

	conf := &restapi.UpdateConf{
		UpdateReq: reqParams,
		ModelObj:  obj,
		//OmitString:       model.OrderOmitUpdate,
		OptModelFunc:     func(db *gorm.DB) *response.ErrorInfo { return optOrderModelForUpdate(db, reqParams, obj) },
		DealAssociations: func(db *gorm.DB) *response.ErrorInfo { return updateOrderProducts(db, reqParams, obj) },
	}

	restapi.Update(ctx, conf)
}

func optOrderModelForUpdate(db *gorm.DB, reqParams *UpdateOrderReq, obj *model.Order) *response.ErrorInfo {
	ei := optCustomerID(db, reqParams, obj)
	if ei != nil {
		return ei
	}

	ei = optSalesmanID(db, reqParams, obj)
	if ei != nil {
		return ei
	}

	ei = optDelivererID(db, reqParams, obj)
	if ei != nil {
		return ei
	}

	utils.OptIntPtr(&obj.PayMode, reqParams.PayMode)
	utils.OptFloat32Ptr(&obj.Paid, reqParams.Paid)
	utils.OptStringPtr(&obj.Comment, reqParams.Comment)
	utils.OptIntPtr(&obj.DeliveryMode, reqParams.DeliveryMode)
	utils.OptTimePtrByUnixTimePtr(&obj.DeliveryAt, reqParams.DeliveryAt)
	utils.OptStringPtr(&obj.DeliveryAddress, reqParams.DeliveryAddress)

	return nil
}

func optCustomerID(db *gorm.DB, reqParams *UpdateOrderReq, obj *model.Order) *response.ErrorInfo {
	if reqParams.CustomerId == nil {
		return nil
	}

	ei := store.EnsureExistByID(db, &model.Customer{}, *reqParams.CustomerId)
	if ei != nil {
		return ei
	}

	obj.CustomerID = utils.GetInt64PtrByStringPtrDefaultNil(reqParams.CustomerId)
	return nil
}

func optSalesmanID(db *gorm.DB, reqParams *UpdateOrderReq, obj *model.Order) *response.ErrorInfo {
	if reqParams.SalesmanId == nil {
		return nil
	}

	ei := store.EnsureExistByID(db, &model.Salesman{}, *reqParams.SalesmanId)
	if ei != nil {
		return ei
	}

	obj.SalesmanID = utils.GetInt64PtrByStringPtrDefaultNil(reqParams.SalesmanId)
	return nil
}

func optDelivererID(db *gorm.DB, reqParams *UpdateOrderReq, obj *model.Order) *response.ErrorInfo {
	if reqParams.DelivererId == nil {
		return nil
	}

	ei := store.EnsureExistByID(db, &model.Deliverer{}, *reqParams.DelivererId)
	if ei != nil {
		return ei
	}

	obj.DelivererID = utils.GetInt64PtrByStringPtrDefaultNil(reqParams.DelivererId)
	return nil
}

func updateOrderProducts(db *gorm.DB, reqParams *UpdateOrderReq, obj *model.Order) *response.ErrorInfo {
	if len(reqParams.Products) == 0 {
		return nil
	}

	list := make([]model.OrderProduct, 0)
	err := db.Where("order_id = ?", obj.ID).Find(&list).Error
	if err != nil {
		return response.NewStorageError(response.StorageErrorCode, err)
	}

	for k, v := range reqParams.Products {
		op := getOrderProductFromModelList(k, list)
		// when not exist, create
		if op == nil {
			ei := addOrderProductForUpdate(db, obj, k, v)
			if ei != nil {
				return ei
			}
			continue
		}

		// when the request product is nil, delete the OrderProduct
		if v == nil {
			err = store.Delete(db, op)
			if err != nil {
				return response.NewStorageError(response.StorageErrorCode, err)
			}
			continue
		}

		//when the request product is not nil, update the OrderProduct
		ei := updateOrderProduct(db, op, v)
		if ei != nil {
			return ei
		}
	}
	return nil
}

func getOrderProductFromModelList(productID string, list []model.OrderProduct) *model.OrderProduct {
	for _, v := range list {
		id, _ := strconv.ParseInt(productID, 10, 64)
		if v.ProductID == id {
			return &v
		}
	}
	return nil
}

func addOrderProductForUpdate(db *gorm.DB, obj *model.Order, productId string, proParam *ProductForUpdate) *response.ErrorInfo {
	if proParam == nil {
		return nil
	}

	pro := model.Product{}
	err := store.GetByID(db, &pro, productId)
	if err != nil {
		return response.NewCommonError(response.NotFoundErrorCode)
	}

	op := model.OrderProduct{
		Model:       store.GenerateModel(),
		OrderID:     obj.ID,
		ProductID:   pro.ID,
		Count:       utils.GetIntFromPtr(proParam.Count),
		BoughtPrice: pro.Price,
		Paid:        utils.GetFloat32FromPtr(proParam.Paid),
		FinalPrice:  utils.GetFloat32FromPtr(proParam.FinalPrice),
		Comment:     utils.GetStringFromPtr(proParam.Comment),
	}

	err = store.Create(db, op)
	if err != nil {
		return response.NewStorageError(response.StorageErrorCode, err)
	}

	return nil
}

func updateOrderProduct(db *gorm.DB, obj *model.OrderProduct, proParam *ProductForUpdate) *response.ErrorInfo {
	if proParam == nil {
		return nil
	}

	utils.OptIntPtr(&obj.Count, proParam.Count)
	utils.OptFloat32Ptr(&obj.Paid, proParam.Paid)
	utils.OptFloat32Ptr(&obj.FinalPrice, proParam.FinalPrice)
	utils.OptStringPtr(&obj.Comment, proParam.Comment)

	err := store.Update(db, obj)
	if err != nil {
		return response.NewStorageError(response.StorageErrorCode, err)
	}

	return nil
}
