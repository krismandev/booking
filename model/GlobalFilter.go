package model

type GlobalQueryFilter struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	OrderDir string `json:"orderDir"`
	OrderBy  string `json:"orderBy"`
}
