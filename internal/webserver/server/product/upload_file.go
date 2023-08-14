package product

import (
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"zq-xu/warehouse-admin/internal/webserver/config"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/utils"
)

type UploadFileToInstanceReq struct {
	FormFileHeader *multipart.FileHeader
}

func UploadFile(ctx *gin.Context) {
	reqParams, ei := newUploadFileReq(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	_, ei = auth.GetAccessControl(ctx, ctx.GetString(auth.AuthUserIDToken))
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	dir, localPath, ei := saveUploadFilesToLocal(ctx, reqParams)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}
	defer os.RemoveAll(dir)

	log.Logger.Infof("Succeed to save the file %s", localPath)
	ctx.JSON(http.StatusCreated, struct{}{})
}

func newUploadFileReq(ctx *gin.Context) (*UploadFileToInstanceReq, *response.ErrorInfo) {
	f, err := ctx.FormFile(types.FileParam)
	if err != nil {
		return nil, response.NewCommonError(response.GetFormFileErrorCode, err)
	}

	return &UploadFileToInstanceReq{
		FormFileHeader: f,
	}, nil
}

func saveUploadFilesToLocal(ctx *gin.Context, reqParams *UploadFileToInstanceReq) (string, string, *response.ErrorInfo) {
	dir := path.Join(config.WebServerCfg.TmpDir, uuid.New().String())
	localPath := path.Join(dir, reqParams.FormFileHeader.Filename)

	err := utils.EnsureDirExist(dir)
	if err != nil {
		return "", "", response.NewCommonError(response.SaveFileErrorCode, err)
	}

	err = ctx.SaveUploadedFile(reqParams.FormFileHeader, localPath)
	if err != nil {
		return "", "", response.NewCommonError(response.SaveFileErrorCode, err)
	}

	return dir, localPath, nil
}
