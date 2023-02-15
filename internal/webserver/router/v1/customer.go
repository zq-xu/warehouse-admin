package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/customer"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	customerGroupPath = "/customers"
	customerPath      = fmt.Sprintf("/:%s", types.IDParam)
	CustomerGroup     = &router.APIGroup{
		RelativePath: customerGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: customer.CreateCustomer},
			{Method: http.MethodDelete, Path: customerPath, Handler: customer.DeleteCustomer},
			{Method: http.MethodGet, Path: customerPath, Handler: customer.GetCustomer},
			{Method: http.MethodGet, Path: "", Handler: customer.ListCustomer},
		},
	}
)

func init() {
	registerAPIGroup(CustomerGroup)
}
