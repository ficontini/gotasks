package data

const (
	minTitleLen       = 5
	minDescriptionLen = 5
)

type Filter interface{}

type CompletionFilter struct {
	Completed bool
}

type NoFilter struct {
}
type AssignationFilter struct {
	AssignedTo string
}
type CompleteFilter struct {
	AssignedTo string
	Completed  bool
}
