package dto

import "go-document-generator/internal/shared/pagination"

type PaginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func MetaFrom(m pagination.Meta) PaginationMeta {
	return PaginationMeta{Page: m.Page, Limit: m.Limit, Total: m.Total}
}
