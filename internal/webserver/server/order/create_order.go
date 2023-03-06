package order

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/store"
	"zq-xu/warehouse-admin/pkg/utils"
)

type CreateOrderReq struct {
	CustomerId string `json:"customerId"`
	SalesmanId string `json:"salesmanId"`

	PayMode int     `json:"payMode"`
	Paid    float32 `json:"paid"`
	Comment string  `json:"comment"`

	DeliveryMode    int             `json:"deliveryMode"`
	DeliveryAt      *utils.UnixTime `json:"deliveryAt"`
	DeliveryAddress string          `json:"deliveryAddress"`
	DelivererId     string          `json:"delivererId"`

	Products []ProductForCreation `json:"products"`

	Customer      model.Customer       `json:"-"`
	Salesman      model.Salesman       `json:"-"`
	Deliverer     model.Deliverer      `json:"-"`
	OrderProducts []model.OrderProduct `json:"-"`
}

func (cr *CreateOrderReq) CustomerID(id int64) {
	cr.CustomerId = fmt.Sprintf("%d", id)
}

func (cr *CreateOrderReq) SalesmanID(id int64) {
	cr.SalesmanId = fmt.Sprintf("%d", id)
}

func (cr *CreateOrderReq) DelivererID(id int64) {
	cr.DelivererId = fmt.Sprintf("%d", id)
}

type ProductForCreation struct {
	ID       string  `json:"id"`
	Count    int     `json:"count"`
	Paid     float32 `json:"paid"`
	Discount int     `json:"discount"`
	Comment  string  `json:"comment"`
}

func CreateOrder(ctx *gin.Context) {
	reqParams, ei := newCreateOrderReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		ei = loadModelForCreation(reqParams, db)
		if ei != nil {
			return ei
		}

		obj, ei := generateOrderModelForCreation(reqParams)
		if ei != nil {
			return ei
		}

		err := db.Omit(model.OrderOmitCreate...).Create(obj).Error
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		log.Logger.Infof("Succeed to create order %d/%s", obj.ID, obj.OrderNo)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func loadModelForCreation(reqParams *CreateOrderReq, db *gorm.DB) *response.ErrorInfo {
	err := store.GetByID(db, &reqParams.Salesman, reqParams.SalesmanId)
	if err != nil {
		return response.NewCommonError(response.NotFoundErrorCode)
	}

	err = store.GetByID(db, &reqParams.Customer, reqParams.CustomerId)
	if err != nil {
		return response.NewCommonError(response.NotFoundErrorCode)
	}

	if reqParams.DelivererId != "" {
		err = store.GetByID(db, &reqParams.Deliverer, reqParams.DelivererId)
		if err != nil {
			return response.NewCommonError(response.NotFoundErrorCode)
		}
	}

	reqParams.OrderProducts = make([]model.OrderProduct, len(reqParams.Products))
	for k, v := range reqParams.Products {
		pro := model.Product{}
		err = store.GetByID(db, &pro, v.ID)
		if err != nil {
			return response.NewCommonError(response.NotFoundErrorCode)
		}

		optOrderProductForCreation(&reqParams.OrderProducts[k], &pro, &v)
	}
	return nil
}

func newCreateOrderReq(ctx *gin.Context) (*CreateOrderReq, *response.ErrorInfo) {
	reqBody := &CreateOrderReq{}
	err := ctx.ShouldBindJSON(reqBody)
	if err != nil {
		return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
	}

	return reqBody, nil
}

func generateOrderModelForCreation(reqParams *CreateOrderReq) (*model.Order, *response.ErrorInfo) {
	t := &model.Order{
		Model:   store.GenerateModel(),
		OrderNo: uuid.New().String(),
		Comment: reqParams.Comment,
	}

	err := copier.Copy(t, reqParams)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	t.CustomerID = reqParams.Customer.ID
	t.SalesmanID = reqParams.Salesman.ID
	t.DelivererID = reqParams.Deliverer.ID
	return t, nil
}

func optOrderProductForCreation(modelObj *model.OrderProduct, proModel *model.Product, proParam *ProductForCreation) {
	modelObj.Model = store.GenerateModel()
	modelObj.ProductID = proModel.ID
	modelObj.BoughtPrice = proModel.Price
	modelObj.Paid = proParam.Paid
	modelObj.Discount = proParam.Discount
	modelObj.Comment = proParam.Comment
}
