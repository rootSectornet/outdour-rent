package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// Params holds pagination parameters.
type Params struct {
	Page    int `json:"page" form:"page"`
	PerPage int `json:"per_page" form:"per_page"`
}

// NewParams creates validated pagination params.
func NewParams(page, perPage int) Params {
	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}
	return Params{Page: page, PerPage: perPage}
}

// Offset calculates the SQL offset.
func (p Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Limit returns the SQL limit.
func (p Params) Limit() int {
	return p.PerPage
}

// Meta holds pagination metadata in the response.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta creates pagination metadata.
func NewMeta(total int64, page, perPage int) *Meta {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

// FromQuery extracts page and per_page from gin query params.
func FromQuery(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	p := NewParams(page, perPage)
	return p.Page, p.PerPage
}
