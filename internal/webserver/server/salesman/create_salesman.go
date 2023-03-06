package salesman

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/store"
)

type CreateSalesmanReq struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Comment     string `json:"comment"`
	InvoiceInfo string `json:"invoiceInfo"`
}

func CreateSalesman(ctx *gin.Context) {
	reqParams, ei := newCreateSalesmanReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		ei = store.EnsureNotExistByName(db, &model.Salesman{}, reqParams.Name)
		if ei != nil {
			return ei
		}

		obj, ei := generateSalesmanModelForCreation(reqParams)
		if ei != nil {
			return ei
		}

		err := store.Create(db, obj)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		log.Logger.Infof("Succeed to create salesman %d/%s", obj.ID, obj.Name)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func newCreateSalesmanReq(ctx *gin.Context) (*CreateSalesmanReq, *response.ErrorInfo) {
	reqBody := &CreateSalesmanReq{}
	err := ctx.ShouldBindJSON(reqBody)
	if err != nil {
		return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
	}

	return reqBody, nil
}

func generateSalesmanModelForCreation(reqParams *CreateSalesmanReq) (*model.Salesman, *response.ErrorInfo) {
	t := &model.Salesman{
		Model: store.GenerateModel(),
	}

	err := copier.Copy(t, reqParams)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return t, nil
}
