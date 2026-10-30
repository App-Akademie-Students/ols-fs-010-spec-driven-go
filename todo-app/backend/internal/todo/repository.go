package todo

import (
	"math"
	"sync"
)

type Repository interface {
	Create(todo Todo) (Todo, error)
	List(filter Filter) ([]Todo, error)
	ListWithFilters(filter Filter, priority PriorityFilter) ([]Todo, error)
	Update(id int64, status *Status, priority *Priority) (Todo, error)
	UpdateStatus(id int64, status Status) (Todo, error)
	Delete(id int64) error
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	todos     map[int64]Todo
	order     []int64
	nextID    int64
	exhausted bool
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		todos:  make(map[int64]Todo),
		nextID: 1,
	}
}

func (repository *InMemoryRepository) Create(todo Todo) (Todo, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if todo.Status != StatusOpen && todo.Status != StatusCompleted {
		return Todo{}, ErrInvalidStatus
	}
	if todo.Priority == "" {
		todo.Priority = PriorityMedium
	}
	if !todo.Priority.IsValid() {
		return Todo{}, ErrInvalidPriority
	}
	if repository.exhausted {
		return Todo{}, ErrIDExhausted
	}

	todo.ID = repository.nextID
	repository.todos[todo.ID] = todo
	repository.order = append(repository.order, todo.ID)
	if repository.nextID == math.MaxInt64 {
		repository.exhausted = true
	} else {
		repository.nextID++
	}
	return todo, nil
}

func (repository *InMemoryRepository) List(filter Filter) ([]Todo, error) {
	return repository.ListWithFilters(filter, PriorityFilterAll)
}

func (repository *InMemoryRepository) ListWithFilters(filter Filter, priorityFilter PriorityFilter) ([]Todo, error) {
	if filter != FilterAll && filter != FilterOpen && filter != FilterCompleted {
		return nil, ErrInvalidFilter
	}
	if priorityFilter != PriorityFilterAll && priorityFilter != PriorityFilterLow &&
		priorityFilter != PriorityFilterMedium && priorityFilter != PriorityFilterHigh {
		return nil, ErrInvalidPriorityFilter
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()

	todos := make([]Todo, 0, len(repository.todos))
	for _, id := range repository.order {
		todo, exists := repository.todos[id]
		if !exists {
			continue
		}
		if todo.Priority == "" {
			todo.Priority = PriorityMedium
		}
		if matchesStatus(filter, todo.Status) && matchesPriority(priorityFilter, todo.Priority) {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (repository *InMemoryRepository) UpdateStatus(id int64, status Status) (Todo, error) {
	return repository.Update(id, &status, nil)
}

func (repository *InMemoryRepository) Update(id int64, status *Status, priority *Priority) (Todo, error) {
	if status == nil && priority == nil {
		return Todo{}, ErrEmptyUpdate
	}
	if status != nil && *status != StatusOpen && *status != StatusCompleted {
		return Todo{}, ErrInvalidStatus
	}
	if priority != nil && !priority.IsValid() {
		return Todo{}, ErrInvalidPriority
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()

	todo, exists := repository.todos[id]
	if !exists {
		return Todo{}, ErrNotFound
	}
	if todo.Priority == "" {
		todo.Priority = PriorityMedium
	}
	if status != nil {
		todo.Status = *status
	}
	if priority != nil {
		todo.Priority = *priority
	}
	repository.todos[id] = todo
	return todo, nil
}

func matchesStatus(filter Filter, status Status) bool {
	return filter == FilterAll || (filter == FilterOpen && status == StatusOpen) ||
		(filter == FilterCompleted && status == StatusCompleted)
}

func matchesPriority(filter PriorityFilter, priority Priority) bool {
	return filter == PriorityFilterAll || string(filter) == string(priority)
}

func (repository *InMemoryRepository) Delete(id int64) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, exists := repository.todos[id]; !exists {
		return ErrNotFound
	}
	delete(repository.todos, id)
	for index, orderedID := range repository.order {
		if orderedID == id {
			copy(repository.order[index:], repository.order[index+1:])
			repository.order = repository.order[:len(repository.order)-1]
			break
		}
	}
	return nil
}
