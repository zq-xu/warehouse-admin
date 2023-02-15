package supplier

import (
	"fmt"
	"net/http"
	"zq-xu/warehouse-admin/pkg/store"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
)

type CreateSupplierReq struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Comment     string `json:"comment"`
	InvoiceInfo string `json:"invoiceInfo"`
}

func CreateSupplier(ctx *gin.Context) {
	reqParams, ei := newCreateSupplierReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		ei = store.EnsureNotExistByName(db, &model.Supplier{}, reqParams.Name)
		if ei != nil {
			return ei
		}

		obj, ei := generateSupplierModelForCreation(reqParams)
		if ei != nil {
			return ei
		}

		err := store.Create(db, obj)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		log.Logger.Infof("Succeed to create supplier %d/%s", obj.ID, obj.Name)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func newCreateSupplierReq(ctx *gin.Context) (*CreateSupplierReq, *response.ErrorInfo) {
	reqBody := &CreateSupplierReq{}
	err := ctx.ShouldBindJSON(reqBody)
	if err != nil {
		return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
	}

	return reqBody, nil
}

func generateSupplierModelForCreation(reqParams *CreateSupplierReq) (*model.Supplier, *response.ErrorInfo) {
	t := &model.Supplier{
		Model: model.GenerateModel(),
	}

	err := copier.Copy(t, reqParams)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return t, nil
}
