package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/product"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	productGroupPath = "/products"
	productPath      = fmt.Sprintf("/:%s", types.IDParam)
	ProductGroup     = &router.APIGroup{
		RelativePath: productGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: product.CreateProduct},
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
}
