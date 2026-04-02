package pagination

import "math"

type Params struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func (p *Params) Defaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p *Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
}

func NewMeta(page, pageSize int, totalCount int64) Meta {
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}
