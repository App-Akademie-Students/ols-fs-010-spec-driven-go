package todo

import (
	"strings"
	"unicode"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (service *Service) Create(title string, priorities ...Priority) (Todo, error) {
	normalizedTitle := strings.TrimFunc(title, unicode.IsSpace)
	if normalizedTitle == "" {
		return Todo{}, ErrInvalidTitle
	}
	if len(priorities) > 1 {
		return Todo{}, ErrInvalidPriority
	}
	priority := PriorityMedium
	if len(priorities) == 1 {
		priority = priorities[0]
	}
	if !priority.IsValid() {
		return Todo{}, ErrInvalidPriority
	}
	return service.repository.Create(Todo{
		Title:    normalizedTitle,
		Status:   StatusOpen,
		Priority: priority,
	})
}

func (service *Service) List(filter Filter) ([]Todo, error) {
	return service.ListWithFilters(filter, PriorityFilterAll)
}

func (service *Service) ListWithFilters(filter Filter, priorityFilter PriorityFilter) ([]Todo, error) {
	if filter == "" {
		filter = FilterAll
	}
	if filter != FilterAll && filter != FilterOpen && filter != FilterCompleted {
		return nil, ErrInvalidFilter
	}
	if priorityFilter == "" {
		priorityFilter = PriorityFilterAll
	}
	if priorityFilter != PriorityFilterAll && priorityFilter != PriorityFilterLow &&
		priorityFilter != PriorityFilterMedium && priorityFilter != PriorityFilterHigh {
		return nil, ErrInvalidPriorityFilter
	}
	return service.repository.ListWithFilters(filter, priorityFilter)
}

func (service *Service) Update(id int64, status *Status, priority *Priority) (Todo, error) {
	if status == nil && priority == nil {
		return Todo{}, ErrEmptyUpdate
	}
	if status != nil && *status != StatusOpen && *status != StatusCompleted {
		return Todo{}, ErrInvalidStatus
	}
	if priority != nil && !priority.IsValid() {
		return Todo{}, ErrInvalidPriority
	}
	return service.repository.Update(id, status, priority)
}

func (service *Service) UpdateStatus(id int64, status Status) (Todo, error) {
	return service.Update(id, &status, nil)
}

func (service *Service) Delete(id int64) error {
	return service.repository.Delete(id)
}
