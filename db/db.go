package db

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
)

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
func (p *Pagination) SetDefaults() {
	if p.Limit <= 0 {
		p.Limit = DEFAULT_LIMIT
	}
	if p.Page <= 0 {
		p.Page = DEFAULT_PAGE
	}
}
func (p *Pagination) getOptions() *options.FindOptions {
	p.SetDefaults()
	opts := &options.FindOptions{}
	opts.SetSkip((p.Page - 1) * p.Limit)
	opts.SetLimit(p.Limit)
	return opts
}
func (p *Pagination) getQuery() string {
	p.SetDefaults()
	offset := (p.Page - 1) * p.Limit
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, offset)
}

type Store struct {
	Task    TaskStore
	User    UserStore
	Project ProjectStore
}
