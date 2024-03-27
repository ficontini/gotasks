package db

import "errors"

const (
	DBNAME        = "gotasks"
	DBURI         = "mongodb://localhost:27017"
	DEFAULT_PAGE  = 1
	DEFAULT_LIMIT = 10
)

// TODO: Review
var ErrorNotFound = errors.New("resource not found")

type Map map[string]any

type Pagination struct {
	Page  int64
	Limit int64
}

// TODO: Review
func (p *Pagination) CheckDefaultPaginationValues() {
	if p.Limit <= 0 {
		p.Limit = DEFAULT_LIMIT
	}
	if p.Page <= 0 {
		p.Page = DEFAULT_PAGE
	}
}

type Store struct {
	Task TaskStore
}
