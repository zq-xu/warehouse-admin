package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/deliverer"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	delivererGroupPath = "/deliverers"
	delivererPath      = fmt.Sprintf("/:%s", types.IDParam)
	DelivererGroup     = &router.APIGroup{
		RelativePath: delivererGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: deliverer.CreateDeliverer},
			{Method: http.MethodDelete, Path: delivererPath, Handler: deliverer.DeleteDeliverer},
			{Method: http.MethodGet, Path: delivererPath, Handler: deliverer.GetDeliverer},
			{Method: http.MethodGet, Path: "", Handler: deliverer.ListDeliverer},
		},
	}
)

func init() {
	registerAPIGroup(DelivererGroup)
}
