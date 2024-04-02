package db

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/types"
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

func SetUpdateMap(update Map) Map {
	return Map{"$set": update}
}
func PushToKey(update Map) Map {
	return Map{"$push": update}
}
func NewMap(key string, value interface{}) Map {
	return Map{key: value}
}

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
func (p *Pagination) getOptions() *options.FindOptions {
	p.CheckDefaultPaginationValues()
	opts := &options.FindOptions{}
	opts.SetSkip((p.Page - 1) * p.Limit)
	opts.SetLimit(p.Limit)
	return opts
}

type Store struct {
	Task    TaskStore
	User    UserStore
	Project ProjectStore
}
type Updater interface {
	Update(context.Context, types.ID, Map) error
}
