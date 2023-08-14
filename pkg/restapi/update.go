package restapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/store"
)

type UpdateConf struct {
	AuthControl

	UpdateReq  interface{}
	ModelObj   interface{}
	OmitString []string

	OptModelFunc     func(db *gorm.DB) *response.ErrorInfo
	DealAssociations func(db *gorm.DB) *response.ErrorInfo
}

func Update(ctx *gin.Context, conf *UpdateConf) {
	ei := conf.AuthControl.Validate(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	id := ctx.Param(IDParam)
	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		err := store.GetByID(db, conf.ModelObj, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewCommonError(response.NotFoundErrorCode)
			}
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		ei := optUpdateModel(ctx, db, conf)
		if ei != nil {
			return nil
		}

		err = store.Update(db, conf.ModelObj, conf.OmitString...)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		if conf.DealAssociations != nil {
			ei = conf.DealAssociations(db)
			if ei != nil {
				return ei
			}
		}

		log.Logger.Infof("Succeed to update obj %+v", id)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

func optUpdateModel(ctx *gin.Context, db *gorm.DB, conf *UpdateConf) *response.ErrorInfo {
	err := ctx.ShouldBindJSON(conf.UpdateReq)
	if err != nil {
		return response.NewCommonError(response.InvalidParametersErrorCode,
			fmt.Sprintf("request body invalid. %v", err))
	}

	return conf.OptModelFunc(db)
}
