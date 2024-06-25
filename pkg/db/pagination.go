package db

import (
	"gorm.io/gorm"
	"math"
)

type Paginate struct {
	Limit      int `json:"limit,omitempty"`
	Page       int `json:"page,omitempty"`
	TotalRows  int `json:"total_rows,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func NewPaginate(limit, page int) *Paginate {
	return &Paginate{
		Limit: limit, Page: page,
	}
}

func (p *Paginate) PaginatedResult(value interface{}, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	if p.Limit == 0 {
		p.Limit = int(totalRows)
	}
	if p.Page == 0 {
		p.Page = 1
	}
	offset := (p.Page - 1) * p.Limit

	p.TotalRows = int(totalRows)
	p.TotalPages = int(math.Ceil(float64(totalRows) / float64(p.Limit)))
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(p.Limit).Offset(offset)
	}
}
