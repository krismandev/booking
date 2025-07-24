package request

import (
	"math"
	"strconv"
)

type GlobalListDataRequest struct {
	Page     string `json:"page"`
	Limit    string `json:"limit"`
	OrderBy  string `json:"orderBy"`
	OrderDir string `json:"orderDir"`
}

func (r *GlobalListDataRequest) CollectMetadata(count int) (perPage int, currentPage int, totalPage int) {
	limit, _ := strconv.Atoi(r.Limit)
	page, _ := strconv.Atoi(r.Page)

	if limit == 0 {
		limit = 10
	}

	if page == 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	return limit, page, totalPages

}
