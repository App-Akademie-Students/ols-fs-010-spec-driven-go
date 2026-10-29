package todo

import (
	"math"
	"sync"
)

type Repository interface {
	Create(todo Todo) (Todo, error)
	List(filter Filter) ([]Todo, error)
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
	if filter != FilterAll && filter != FilterOpen && filter != FilterCompleted {
		return nil, ErrInvalidFilter
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()

	todos := make([]Todo, 0, len(repository.todos))
	for _, id := range repository.order {
		todo, exists := repository.todos[id]
		if !exists {
			continue
		}
		if filter == FilterAll || (filter == FilterOpen && todo.Status == StatusOpen) ||
			(filter == FilterCompleted && todo.Status == StatusCompleted) {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (repository *InMemoryRepository) UpdateStatus(id int64, status Status) (Todo, error) {
	if status != StatusOpen && status != StatusCompleted {
		return Todo{}, ErrInvalidStatus
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()

	todo, exists := repository.todos[id]
	if !exists {
		return Todo{}, ErrNotFound
	}
	todo.Status = status
	repository.todos[id] = todo
	return todo, nil
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
