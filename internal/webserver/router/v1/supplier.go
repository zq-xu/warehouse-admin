package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/supplier"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	supplierGroupPath = "/suppliers"
	supplierPath      = fmt.Sprintf("/:%s", types.IDParam)
	SupplierGroup     = &router.APIGroup{
		RelativePath: supplierGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: supplier.CreateSupplier},
			{Method: http.MethodDelete, Path: supplierPath, Handler: supplier.DeleteSupplier},
			{Method: http.MethodGet, Path: supplierPath, Handler: supplier.GetSupplier},
			{Method: http.MethodGet, Path: "", Handler: supplier.ListSupplier},
		},
	}
)

func init() {
	registerAPIGroup(SupplierGroup)
}
