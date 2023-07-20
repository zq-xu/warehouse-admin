package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/category"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	categoryGroupPath = "/categories"
	categoryPath      = fmt.Sprintf("/:%s", types.IDParam)
	CategoryGroup     = &router.APIGroup{
		RelativePath: categoryGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: category.CreateCategory},
			{Method: http.MethodPut, Path: categoryPath, Handler: category.UpdateCategory},
			{Method: http.MethodDelete, Path: categoryPath, Handler: category.DeleteCategory},
			{Method: http.MethodGet, Path: categoryPath, Handler: category.GetCategory},
			{Method: http.MethodGet, Path: "", Handler: category.ListCategory},
		},
	}
)

func init() {
	registerAPIGroup(CategoryGroup)
}
