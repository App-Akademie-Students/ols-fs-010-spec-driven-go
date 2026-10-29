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
	if created.ID < 1 {
		t.Errorf("Create() ID = %d, want a positive ID", created.ID)
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
