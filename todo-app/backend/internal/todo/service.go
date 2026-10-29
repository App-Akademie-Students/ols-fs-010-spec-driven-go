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

func (service *Service) Create(title string) (Todo, error) {
	normalizedTitle := strings.TrimFunc(title, unicode.IsSpace)
	if normalizedTitle == "" {
		return Todo{}, ErrInvalidTitle
	}
	return service.repository.Create(Todo{
		Title:  normalizedTitle,
		Status: StatusOpen,
	})
}

func (service *Service) List(filter Filter) ([]Todo, error) {
	if filter == "" {
		filter = FilterAll
	}
	if filter != FilterAll && filter != FilterOpen && filter != FilterCompleted {
		return nil, ErrInvalidFilter
	}
	return service.repository.List(filter)
}

func (service *Service) UpdateStatus(id int64, status Status) (Todo, error) {
	if status != StatusOpen && status != StatusCompleted {
		return Todo{}, ErrInvalidStatus
	}
	return service.repository.UpdateStatus(id, status)
}

func (service *Service) Delete(id int64) error {
	return service.repository.Delete(id)
}
