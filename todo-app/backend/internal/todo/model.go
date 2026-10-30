package todo

import "errors"

type Status string

const (
	StatusOpen      Status = "open"
	StatusCompleted Status = "completed"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (priority Priority) IsValid() bool {
	return priority == PriorityLow || priority == PriorityMedium || priority == PriorityHigh
}

type Filter string

const (
	FilterAll       Filter = "all"
	FilterOpen      Filter = "open"
	FilterCompleted Filter = "completed"
)

type PriorityFilter string

const (
	PriorityFilterAll    PriorityFilter = "all"
	PriorityFilterLow    PriorityFilter = "low"
	PriorityFilterMedium PriorityFilter = "medium"
	PriorityFilterHigh   PriorityFilter = "high"
)

type Todo struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	Status   Status   `json:"status"`
	Priority Priority `json:"priority"`
}

var (
	ErrInvalidTitle          = errors.New("title is required")
	ErrInvalidStatus         = errors.New("invalid todo status")
	ErrInvalidPriority       = errors.New("invalid todo priority")
	ErrInvalidFilter         = errors.New("invalid status filter")
	ErrInvalidPriorityFilter = errors.New("invalid priority filter")
	ErrEmptyUpdate           = errors.New("todo update has no fields")
	ErrNotFound              = errors.New("todo not found")
	ErrIDExhausted           = errors.New("todo ID space exhausted")
)
