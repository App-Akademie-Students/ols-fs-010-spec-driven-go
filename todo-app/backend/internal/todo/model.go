package todo

import "errors"

type Status string

const (
	StatusOpen      Status = "open"
	StatusCompleted Status = "completed"
)

type Filter string

const (
	FilterAll       Filter = "all"
	FilterOpen      Filter = "open"
	FilterCompleted Filter = "completed"
)

type Todo struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status Status `json:"status"`
}

var (
	ErrInvalidTitle  = errors.New("title is required")
	ErrInvalidStatus = errors.New("invalid todo status")
	ErrInvalidFilter = errors.New("invalid status filter")
	ErrNotFound      = errors.New("todo not found")
	ErrIDExhausted   = errors.New("todo ID space exhausted")
)
