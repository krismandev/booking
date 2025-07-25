package model

type GlobalQueryFilter struct {
	Page     string `json:"page"`
	Limit    string `json:"limit"`
	OrderDir string `json:"orderDir"`
	OrderBy  string `json:"orderBy"`
}
