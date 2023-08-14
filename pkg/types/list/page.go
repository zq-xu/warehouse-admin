package list

const (
	defaultPageSize = 10

	PageNumQuery  = "page_num"
	PageSizeQuery = "page_size"
)

type PageResponse struct {
	PageInfo `json:",inline"`

	Count int64       `json:"count"`
	Items interface{} `json:"items"`
}

type PageInfo struct {
	PageQueries `json:",inline"`
	PageCount   int64 `json:"page_total"`
}

type PageQueries struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

func NewPageResponse(count int64, pi PageQueries, items interface{}) *PageResponse {
	pr := &PageResponse{
		PageInfo: PageInfo{
			PageQueries: pi,
		},

		Count: count,
		Items: items,
	}

	if pr.PageSize != 0 {
		pr.PageCount = pr.Count / int64(pr.PageSize)
	}

	if pr.PageCount*int64(pr.PageSize) < pr.Count {
		pr.PageCount++
	}

	return pr
}

func (p *PageInfo) Revise() {
	if p.PageSize == 0 {
		p.PageSize = defaultPageSize
	}

	if p.PageNum <= 0 {
		p.PageNum = 1
	}
}
