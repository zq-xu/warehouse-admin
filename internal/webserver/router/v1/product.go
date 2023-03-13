package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/internal/webserver/server/product"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
	"zq-xu/warehouse-admin/pkg/router/auditlog"
)

var (
	productGroupPath = "/products"
	productPath      = fmt.Sprintf("/:%s", types.IDParam)
	ProductGroup     = &router.APIGroup{
		RelativePath: productGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: product.CreateProduct},
			{Method: http.MethodPut, Path: productPath, Handler: product.UpdateProduct},
			{Method: http.MethodDelete, Path: productPath, Handler: product.DeleteProduct},
			{Method: http.MethodGet, Path: productPath, Handler: product.GetProduct},
			{Method: http.MethodGet, Path: "", Handler: product.ListProduct},
			{Method: http.MethodPost, Path: "/upload", Handler: product.UploadFile},
			{Method: http.MethodPost, Path: "/stockin", Handler: product.StockInProduct},
		},
	}
)

func init() {
	registerAPIGroup(ProductGroup)
	auditlog.Middleware.RegisterBodyGenerateFn(http.MethodPost, VersionV1+productGroupPath, func(ctx *gin.Context) []byte {
		r, _ := product.GenerateBaseReq(ctx)
		bs, _ := json.Marshal(r)
		return bs
	})
}
