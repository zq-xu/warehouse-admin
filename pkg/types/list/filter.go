package list

import "zq-xu/warehouse-admin/pkg/constant"

type filter struct {
	filterParamSet map[string]constant.FilterRule
}

type FilterObj interface {
	GetFilterString(f string) string
}

func NewFilter(fs map[string]constant.FilterRule) *filter {
	tmp := fs
	if tmp == nil {
		tmp = make(map[string]constant.FilterRule)
	}

	return &filter{
		filterParamSet: tmp,
	}
}

func (f *filter) IsMatch(fps constant.FilterQueries, obj FilterObj) bool {
	for k, v := range fps {
		s := obj.GetFilterString(k)

		rule, ok := f.filterParamSet[k]
		if !ok {
			continue
		}

		if !rule.IsMatch(s, v) {
			return false
		}
	}

	return true
}
