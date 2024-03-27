package api

import (
	"reflect"
)

type ResourceResponse struct {
	Data    any   `json:"data"`
	Results int   `json:"results"`
	Page    int64 `json:"page"`
}

func NewResourceResponse(data any, results int, page int64) ResourceResponse {
	//TODO: Review
	if reflect.ValueOf(data).IsNil() {
		data = []any{}
	}
	return ResourceResponse{
		Data:    data,
		Results: results,
		Page:    page,
	}
}
