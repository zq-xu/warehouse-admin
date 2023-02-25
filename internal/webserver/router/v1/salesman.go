package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/salesman"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	salesmanGroupPath = "/salesmen"
	salesmanPath      = fmt.Sprintf("/:%s", types.IDParam)
	SalesmanGroup     = &router.APIGroup{
		RelativePath: salesmanGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: salesman.CreateSalesman},
			{Method: http.MethodPut, Path: salesmanPath, Handler: salesman.UpdateSalesman},
			{Method: http.MethodDelete, Path: salesmanPath, Handler: salesman.DeleteSalesman},
			{Method: http.MethodGet, Path: salesmanPath, Handler: salesman.GetSalesman},
			{Method: http.MethodGet, Path: "", Handler: salesman.ListSalesman},
		},
	}
)

func init() {
	registerAPIGroup(SalesmanGroup)
}
