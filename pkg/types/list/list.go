package list

import (
	"zq-xu/warehouse-admin/pkg/constant"
)

type Params struct {
	PageQueries   PageQueries
	FilterQueries constant.FilterQueries
	SortQuery     string
}
