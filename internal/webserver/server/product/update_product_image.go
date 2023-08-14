package product

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/store"
)

const (
	IdFormKey = "id"
)

type UpdateProductImageReq struct {
	Id string `json:"-"`

	userId string `json:"-"`

	obj model.Product `json:"-"`
}

func UpdateProductImage(ctx *gin.Context) {
	_, ei := auth.GetAccessControl(ctx, ctx.GetString(auth.AuthUserIDToken))
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ei = store.DoDBTransaction(store.DB(ctx), func(db *gorm.DB) *response.ErrorInfo {
		reqParams, ei := newUpdateProductImageReq(ctx, db)
		if ei != nil {
			return ei
		}

		err := store.Update(db, &reqParams.obj)
		if err != nil {
			return response.NewStorageError(response.StorageErrorCode, err)
		}

		log.Logger.Infof("Succeed to update image for the product  %d/%s", reqParams.obj.ID, reqParams.obj.Name)
		return nil
	})
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	ctx.JSON(http.StatusCreated, struct{}{})
}

// GenerateUpdateProductImageBaseReq
// for audit log
func GenerateUpdateProductImageBaseReq(ctx *gin.Context) *UpdateProductImageReq {
	reqBody := &UpdateProductImageReq{}
	reqBody.Id = ctx.PostForm(IdFormKey)
	reqBody.userId = auth.GetAuthUserId(ctx)
	return reqBody
}

func newUpdateProductImageReq(ctx *gin.Context, db *gorm.DB) (*UpdateProductImageReq, *response.ErrorInfo) {
	reqBody := GenerateUpdateProductImageBaseReq(ctx)

	// check if the product exists
	err := store.GetByID(db, &reqBody.obj, reqBody.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCommonError(response.NotFoundErrorCode)
		}
		return nil, response.NewStorageError(response.StorageErrorCode, err)
	}

	cr := &CreateProductReq{modelObj: reqBody.obj.Model, userId: reqBody.userId}

	isFileExist, ei := uploadImageToS3(ctx, cr)
	if ei != nil {
		return nil, ei
	}
	// the image is not uploaded
	if !isFileExist {
		return nil, response.NewCommonError(response.GetFormFileErrorCode, http.ErrMissingFile)
	}

	ei = uploadThumbnailToS3(ctx, cr)
	if ei != nil {
		return nil, ei
	}

	reqBody.obj.Image = cr.Image
	reqBody.obj.Thumbnail = cr.Thumbnail
	return reqBody, ei
}
