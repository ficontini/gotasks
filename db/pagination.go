package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DEFAULT_PAGE  = 1
	DEFAULT_LIMIT = 10
)

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
func (p *Pagination) generatePagination() (int64, int64) {
	p.SetDefaults()
	skip := (p.Page - 1) * p.Limit
	limit := p.Limit
	return skip, limit
}
func (p *Pagination) generatePaginationForDynamoDB() (int32, *int32) {
	p.SetDefaults()
	offset := int32((p.Page - 1) * p.Limit)
	limit := int32(p.Limit)
	return offset, &limit
}

func (p *Pagination) getOptions() *options.FindOptions {
	skip, limit := p.generatePagination()
	opts := &options.FindOptions{}
	opts.SetSkip(skip)
	opts.SetLimit(limit)
	return opts
}
func (p *Pagination) getQuery() string {
	skip, limit := p.generatePagination()
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, skip)
}