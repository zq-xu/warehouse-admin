package product

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/internal/webserver/config"
	"zq-xu/warehouse-admin/internal/webserver/model"
	"zq-xu/warehouse-admin/pkg/awsapi"
	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/store"
	"zq-xu/warehouse-admin/pkg/utils"
)

const (
	NameFormKey           = "name"
	PriceFormKey          = "price"
	ImageFormKey          = "image"
	StorageAddressFormKey = "storageAddress"
	CommentFormKey        = "comment"

	thumbnailWidth  = 100
	thumbnailHeight = 100

	productImageSubDir = "products/images"
)

var productImageFormatSuffix = []string{".jpg", ".jpeg", ".png", ".svg"}

type CreateProductReq struct {
	Name           string  `json:"name"`
	Price          float32 `json:"price"`
	StorageAddress string  `json:"storageAddress"`
	Comment        string  `json:"comment"`

	Image     string `json:"-"`
	Thumbnail string `json:"-"`

	userId            string
	imageFormatSuffix string
	modelObj          store.Model
}

func CreateProduct(ctx *gin.Context) {
	reqParams, ei := newCreateProductReq(ctx)
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

// GenerateBaseReq
// for audit log
func GenerateBaseReq(ctx *gin.Context) (*CreateProductReq, *response.ErrorInfo) {
	reqBody := &CreateProductReq{}

	priceStr := ctx.PostForm(PriceFormKey)
	f32, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		return reqBody, response.NewCommonError(response.InvalidParametersErrorCode, fmt.Sprintf("request price invalid. %v", err))
	}

	reqBody.Price = float32(f32)

	reqBody.Name = ctx.PostForm(NameFormKey)
	reqBody.StorageAddress = ctx.PostForm(StorageAddressFormKey)
	reqBody.Comment = ctx.PostForm(CommentFormKey)

	reqBody.userId = auth.GetAuthUserId(ctx)
	return reqBody, nil
}

func newCreateProductReq(ctx *gin.Context) (*CreateProductReq, *response.ErrorInfo) {
	reqBody, ei := GenerateBaseReq(ctx)
	if ei != nil {
		return nil, ei
	}
	reqBody.modelObj = store.GenerateModel()

	isFileExist, ei := uploadImageToS3(ctx, reqBody)
	if !isFileExist {
		return reqBody, ei
	}
	if ei != nil {
		return nil, ei
	}

	ei = uploadThumbnailToS3(ctx, reqBody)
	if ei != nil {
		return nil, ei
	}

	return reqBody, ei
}

func getFileFormatSuffix(filename string) string {
	index := strings.LastIndex(filename, ".")

	if index < 0 {
		return ""
	}

	return filename[index:]
}

func openImageFile(ctx *gin.Context) (string, multipart.File, *response.ErrorInfo) {
	fh, err := ctx.FormFile(ImageFormKey)
	if err != nil {
		return "", nil, response.NewCommonError(response.GetFormFileErrorCode, err)
	}

	f, err := fh.Open()
	if err != nil {
		return "", nil, response.NewCommonError(response.GetFormFileErrorCode, err.Error())
	}

	return fh.Filename, f, nil
}

func uploadImageToS3(ctx *gin.Context, reqParams *CreateProductReq) (bool, *response.ErrorInfo) {
	fh, err := ctx.FormFile(ImageFormKey)
	if err != nil {
		if err == http.ErrMissingFile {
			return false, nil
		}
		return true, response.NewCommonError(response.GetFormFileErrorCode, err)
	}

	f, err := fh.Open()
	if err != nil {
		return true, response.NewCommonError(response.GetFormFileErrorCode, err)
	}
	defer f.Close()

	reqParams.imageFormatSuffix = getFileFormatSuffix(fh.Filename)
	if !utils.ContainString(productImageFormatSuffix, reqParams.imageFormatSuffix) {
		return true, response.NewCommonError(response.InvalidImageFormatErrorCode)
	}

	imgS3Path := awsapi.GenerateS3BucketPath(reqParams.userId, productImageSubDir,
		fmt.Sprintf("%d%s", reqParams.modelObj.ID, reqParams.imageFormatSuffix))

	op, err := awsapi.S3Client.UploadFileByReader(f, awsapi.S3Cfg.Bucket, imgS3Path)
	if err != nil {
		return true, response.NewCommonError(response.UploadFileToS3ErrorCode, err.Error())
	}

	reqParams.Image = op.Location
	return true, nil
}

func uploadThumbnailToS3(ctx *gin.Context, reqParams *CreateProductReq) *response.ErrorInfo {
	thumbnailName := fmt.Sprintf("%d_tmp%s", reqParams.modelObj.ID, reqParams.imageFormatSuffix)
	thumbnailPath, ei := generateThumbnail(ctx, thumbnailName)
	if ei != nil {
		return ei
	}
	defer os.Remove(thumbnailPath)

	thumbnailS3Path := awsapi.GenerateS3BucketPath(reqParams.userId, productImageSubDir, thumbnailName)
	op, err := awsapi.S3Client.UploadFile(thumbnailPath, awsapi.S3Cfg.Bucket, thumbnailS3Path)
	if err != nil {
		return response.NewCommonError(response.UploadFileToS3ErrorCode, err.Error())
	}

	reqParams.Thumbnail = op.Location
	return nil
}

func generateThumbnail(ctx *gin.Context, thumbnailName string) (string, *response.ErrorInfo) {
	_, f, ei := openImageFile(ctx)
	if ei != nil {
		return "", ei
	}
	defer f.Close()

	localPath := path.Join(config.WebServerCfg.TmpDir, thumbnailName)
	return localPath, generateAndSaveThumbnail(f, localPath)
}

func generateAndSaveThumbnail(r io.Reader, localPath string) *response.ErrorInfo {
	img, err := imaging.Decode(r)
	if err != nil {
		return response.NewCommonError(response.ResizeFileErrorCode, err)
	}

	img = imaging.Resize(img, thumbnailWidth, thumbnailHeight, imaging.Lanczos)
	err = imaging.Save(img, localPath)
	if err != nil {
		return response.NewCommonError(response.SaveFileErrorCode, err)
	}

	return nil
}

func generateProductModelForCreation(reqParams *CreateProductReq) (*model.Product, *response.ErrorInfo) {
	t := &model.Product{
		Model: reqParams.modelObj,
	}

	err := copier.Copy(t, reqParams)
	if err != nil {
		return nil, response.NewCommonError(response.GenerateModelErrorCode, err.Error())
	}

	return t, nil
}
