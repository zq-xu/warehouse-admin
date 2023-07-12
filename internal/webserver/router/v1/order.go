package v1

import (
	"fmt"
	"net/http"

	"zq-xu/warehouse-admin/internal/webserver/server/order"
	"zq-xu/warehouse-admin/internal/webserver/types"
	"zq-xu/warehouse-admin/pkg/router"
)

var (
	orderGroupPath = "/orders"
	orderPath      = fmt.Sprintf("/:%s", types.IDParam)
	OrderGroup     = &router.APIGroup{
		RelativePath: orderGroupPath,
		APIs: []*router.API{
			{Method: http.MethodPost, Path: "", Handler: order.CreateOrder},
			{Method: http.MethodPut, Path: orderPath, Handler: order.UpdateOrder},
			{Method: http.MethodDelete, Path: orderPath, Handler: order.DeleteOrder},
			{Method: http.MethodGet, Path: orderPath, Handler: order.GetOrder},
			{Method: http.MethodGet, Path: "", Handler: order.ListOrder},
			{Method: http.MethodPost, Path: "/stockout", Handler: order.StockOut},
		},
	}
)

func init() {
	registerAPIGroup(OrderGroup)
}
