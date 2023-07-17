package list

import (
	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/pkg/restapi/response"
)

const (
	SortByQuery    = "sort_by"
	defaultSortStr = "updated_at,desc"
)

type Params struct {
	PageInfo  *PageInfo
	Queries   Queries
	SortQuery string
}

func GetListParams(ctx *gin.Context) (*Params, *response.ErrorInfo) {
	pi, ei := GetPageInfo(ctx)
	if ei != nil {
		return nil, ei
	}

	return &Params{
		PageInfo:  pi,
		Queries:   GetQueries(ctx),
		SortQuery: ctx.DefaultQuery(SortByQuery, defaultSortStr),
	}, nil
}
