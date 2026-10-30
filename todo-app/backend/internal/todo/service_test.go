package todo

import (
	"errors"
	"testing"
)

func TestServiceCreateNormalizesTitleAndStartsOpen(t *testing.T) {
	service := NewService(NewInMemoryRepository())

	created, err := service.Create("\u00a0\u2003Prepare demo\u00a0 ")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Title != "Prepare demo" {
		t.Errorf("Create() title = %q, want %q", created.Title, "Prepare demo")
	}
	if created.Status != StatusOpen {
		t.Errorf("Create() status = %q, want %q", created.Status, StatusOpen)
	}
	if created.Priority != PriorityMedium {
		t.Errorf("Create() priority = %q, want %q", created.Priority, PriorityMedium)
	}
	if created.ID < 1 {
		t.Errorf("Create() ID = %d, want a positive ID", created.ID)
	}
}

func TestServiceCreatesAndUpdatesPriorityWithoutChangingOtherFields(t *testing.T) {
	service := NewService(NewInMemoryRepository())
	priorities := []Priority{PriorityLow, PriorityMedium, PriorityHigh}
	createdTodos := make([]Todo, 0, len(priorities))
	for _, priority := range priorities {
		created, err := service.Create("Priority task", priority)
		if err != nil {
			t.Fatalf("Create(%q) error = %v", priority, err)
		}
		if created.Priority != priority {
			t.Errorf("Create(%q) priority = %q", priority, created.Priority)
		}
		createdTodos = append(createdTodos, created)
	}

	newPriority := PriorityHigh
	updated, err := service.Update(createdTodos[0].ID, nil, &newPriority)
	if err != nil {
		t.Fatalf("Update(priority only) error = %v", err)
	}
	if updated.Priority != PriorityHigh || updated.Title != createdTodos[0].Title || updated.Status != createdTodos[0].Status {
		t.Errorf("Update(priority only) = %#v; want priority changed and other fields preserved", updated)
	}

	invalidPriority := Priority("urgent")
	if _, err := service.Create("Invalid priority", invalidPriority); !errors.Is(err, ErrInvalidPriority) {
		t.Errorf("Create(invalid priority) error = %v, want ErrInvalidPriority", err)
	}
	if _, err := service.Create("Empty priority", Priority("")); !errors.Is(err, ErrInvalidPriority) {
		t.Errorf("Create(empty priority) error = %v, want ErrInvalidPriority", err)
	}
	if _, err := service.Update(createdTodos[1].ID, nil, &invalidPriority); !errors.Is(err, ErrInvalidPriority) {
		t.Errorf("Update(invalid priority) error = %v, want ErrInvalidPriority", err)
	}
}

func TestServiceRejectsBlankTitleWithoutCreatingTodo(t *testing.T) {
	repository := NewInMemoryRepository()
	service := NewService(repository)

	_, err := service.Create("\u00a0\u2003\t")
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("Create() error = %v, want ErrInvalidTitle", err)
	}
	todos, err := service.List(FilterAll)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("List() returned %d todos after rejected create, want 0", len(todos))
	}
}

func TestServiceAllowsDuplicateTitlesAndPreservesCreationOrder(t *testing.T) {
	service := NewService(NewInMemoryRepository())
	first, err := service.Create("Same title")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Create("Same title")
	if err != nil {
		t.Fatal(err)
	}

	todos, err := service.List(FilterAll)
	if err != nil {
		t.Fatal(err)
	}
	if len(todos) != 2 || todos[0].ID != first.ID || todos[1].ID != second.ID {
		t.Fatalf("List() = %#v, want both todos in creation order", todos)
	}
	if first.ID == second.ID {
		t.Fatalf("duplicate titles received the same ID %d", first.ID)
	}
	defaultFilter, err := service.List("")
	if err != nil || len(defaultFilter) != 2 {
		t.Fatalf("List(omitted filter) = %#v, %v; want both todos", defaultFilter, err)
	}
}

func TestServiceTransitionsDeleteAndFilters(t *testing.T) {
	service := NewService(NewInMemoryRepository())
	first, err := service.Create("First")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Create("Second")
	if err != nil {
		t.Fatal(err)
	}

	updated, err := service.UpdateStatus(first.ID, StatusCompleted)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != StatusCompleted {
		t.Fatalf("UpdateStatus() status = %q, want %q", updated.Status, StatusCompleted)
	}
	open, err := service.List(FilterOpen)
	if err != nil || len(open) != 1 || open[0].ID != second.ID {
		t.Fatalf("List(open) = %#v, %v; want second todo", open, err)
	}
	completed, err := service.List(FilterCompleted)
	if err != nil || len(completed) != 1 || completed[0].ID != first.ID {
		t.Fatalf("List(completed) = %#v, %v; want first todo", completed, err)
	}

	reopened, err := service.UpdateStatus(first.ID, StatusOpen)
	if err != nil {
		t.Fatal(err)
	}
	if reopened.Status != StatusOpen {
		t.Fatalf("reopened status = %q, want %q", reopened.Status, StatusOpen)
	}
	if err := service.Delete(second.ID); err != nil {
		t.Fatal(err)
	}
	all, err := service.List(FilterAll)
	if err != nil || len(all) != 1 || all[0].ID != first.ID {
		t.Fatalf("List(all) after delete = %#v, %v; want first todo only", all, err)
	}
}

func TestServiceFiltersByStatusAndPriority(t *testing.T) {
	service := NewService(NewInMemoryRepository())
	low, err := service.Create("low open", PriorityLow)
	if err != nil {
		t.Fatal(err)
	}
	highOpen, err := service.Create("high open", PriorityHigh)
	if err != nil {
		t.Fatal(err)
	}
	highCompleted, err := service.Create("high completed", PriorityHigh)
	if err != nil {
		t.Fatal(err)
	}
	completed := StatusCompleted
	if _, err := service.Update(highCompleted.ID, &completed, nil); err != nil {
		t.Fatal(err)
	}

	filtered, err := service.ListWithFilters(FilterOpen, PriorityFilterHigh)
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].ID != highOpen.ID {
		t.Fatalf("ListWithFilters(open, high) = %#v, want only high open todo", filtered)
	}
	filtered, err = service.ListWithFilters(FilterCompleted, PriorityFilterHigh)
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].ID != highCompleted.ID {
		t.Fatalf("ListWithFilters(completed, high) = %#v, want only high completed todo", filtered)
	}
	filtered, err = service.ListWithFilters(FilterAll, PriorityFilterMedium)
	if err != nil || len(filtered) != 0 {
		t.Fatalf("ListWithFilters(all, medium) = %#v, %v; want empty results", filtered, err)
	}
	if _, err := service.ListWithFilters(FilterAll, PriorityFilter("urgent")); !errors.Is(err, ErrInvalidPriorityFilter) {
		t.Errorf("ListWithFilters(invalid priority) error = %v, want ErrInvalidPriorityFilter", err)
	}
	if low.Priority != PriorityLow {
		t.Errorf("low todo priority = %q, want %q", low.Priority, PriorityLow)
	}
}

func TestServiceRejectsInvalidStatusesFiltersAndMissingIDs(t *testing.T) {
	service := NewService(NewInMemoryRepository())
	created, err := service.Create("Existing")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.UpdateStatus(created.ID, "archived"); !errors.Is(err, ErrInvalidStatus) {
		t.Errorf("UpdateStatus() error = %v, want ErrInvalidStatus", err)
	}
	if _, err := service.List("archived"); !errors.Is(err, ErrInvalidFilter) {
		t.Errorf("List() error = %v, want ErrInvalidFilter", err)
	}
	if _, err := service.UpdateStatus(created.ID+1, StatusOpen); !errors.Is(err, ErrNotFound) {
		t.Errorf("UpdateStatus(missing) error = %v, want ErrNotFound", err)
	}
	if err := service.Delete(created.ID + 1); !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete(missing) error = %v, want ErrNotFound", err)
	}
}
