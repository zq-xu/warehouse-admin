package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/productlot"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	productLotGroupPath = "/productlots"
	productLotPath      = fmt.Sprintf("/:%s", types.IDParam)
	ProductLotGroup     = &router.APIGroup{
		RelativePath: productLotGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: productlot.CreateProductLot},
			{Method: http.MethodDelete, Path: productLotPath, Handler: productlot.DeleteProductLot},
			{Method: http.MethodGet, Path: productLotPath, Handler: productlot.GetProductLot},
			{Method: http.MethodGet, Path: "", Handler: productlot.ListProductLot},
		},
	}
)

func init() {
	registerAPIGroup(ProductLotGroup)
}
