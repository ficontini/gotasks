package db

import "errors"

const (
	DBNAME = "gotasks"
	DBURI  = "mongodb://localhost:27017"
)

// TODO: Review
var ErrorNotFound = errors.New("resource not found")

type Map map[string]any

type Store struct {
	Task TaskStore
}
