package request

import (
	"math"
)

type GlobalListDataRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	OrderBy  string `json:"orderBy"`
	OrderDir string `json:"orderDir"`
}

func (r *GlobalListDataRequest) CollectMetadata(count int) (perPage int, currentPage int, totalPage int) {
	limit := r.Limit
	page := r.Page

	if limit == 0 {
		limit = 10
	}

	if page == 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	return limit, page, totalPages

}
