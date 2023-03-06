package product

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

type CreateProductReq struct {
	Name           string  `json:"name"`
	Image          string  `json:"image"`
	Price          float32 `json:"price"`
	StorageAddress string  `json:"storageAddress"`
	Comment        string  `json:"comment"`
}

func CreateProduct(ctx *gin.Context) {
	reqParams, ei := newCreateProductReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		ei = store.EnsureNotExistByName(db, &model.Product{}, reqParams.Name)
		if ei != nil {
			return ei
		}

		obj, ei := generateProductModelForCreation(reqParams)
		if ei != nil {
			return ei
		}

		err := store.Create(db, obj)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		log.Logger.Infof("Succeed to create product %d/%s", obj.ID, obj.Name)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func newCreateProductReq(ctx *gin.Context) (*CreateProductReq, *response.ErrorInfo) {
	reqBody := &CreateProductReq{}
	err := ctx.ShouldBindJSON(reqBody)
	if err != nil {
		return nil, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request body invalid. %v", err))
	}

	return reqBody, nil
}

func generateProductModelForCreation(reqParams *CreateProductReq) (*model.Product, *response.ErrorInfo) {
	t := &model.Product{
		Model: store.GenerateModel(),
	}

	err := copier.Copy(t, reqParams)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return t, nil
}
