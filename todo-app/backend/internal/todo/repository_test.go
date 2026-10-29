package todo

import (
	"sync"
	"testing"
)

func TestRepositoryReturnsDetachedValuesAndFiltersInCreationOrder(t *testing.T) {
	repository := NewInMemoryRepository()
	first, err := repository.Create(Todo{Title: "first", Status: StatusOpen})
	if err != nil {
		t.Fatal(err)
	}
	second, err := repository.Create(Todo{Title: "second", Status: StatusCompleted})
	if err != nil {
		t.Fatal(err)
	}

	first.Title = "mutated by caller"
	all, err := repository.List(FilterAll)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].ID != first.ID || all[1].ID != second.ID || all[0].Title != "first" {
		t.Fatalf("List(all) = %#v, want detached values in creation order", all)
	}
	all[0].Title = "mutated snapshot"
	open, err := repository.List(FilterOpen)
	if err != nil {
		t.Fatal(err)
	}
	if len(open) != 1 || open[0].Title != "first" {
		t.Fatalf("List(open) = %#v, want original first todo", open)
	}
}

func TestRepositoryRemovesDeletedIDsFromCreationOrder(t *testing.T) {
	repository := NewInMemoryRepository()
	first, err := repository.Create(Todo{Title: "first", Status: StatusOpen})
	if err != nil {
		t.Fatal(err)
	}
	second, err := repository.Create(Todo{Title: "second", Status: StatusOpen})
	if err != nil {
		t.Fatal(err)
	}
	third, err := repository.Create(Todo{Title: "third", Status: StatusOpen})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Delete(second.ID); err != nil {
		t.Fatal(err)
	}
	fourth, err := repository.Create(Todo{Title: "fourth", Status: StatusOpen})
	if err != nil {
		t.Fatal(err)
	}

	todos, err := repository.List(FilterAll)
	if err != nil {
		t.Fatal(err)
	}
	want := []int64{first.ID, third.ID, fourth.ID}
	if len(todos) != len(want) {
		t.Fatalf("List() returned %d todos, want %d", len(todos), len(want))
	}
	for index, todo := range todos {
		if todo.ID != want[index] {
			t.Errorf("List()[%d].ID = %d, want %d", index, todo.ID, want[index])
		}
	}
}

func TestRepositoryConcurrentOperations(t *testing.T) {
	repository := NewInMemoryRepository()
	const workers = 8
	const perWorker = 80

	start := make(chan struct{})
	finished := make(chan struct{})
	var workersDone sync.WaitGroup
	var listDone sync.WaitGroup
	createdIDs := make(chan int64, workers*perWorker)
	workersDone.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			defer workersDone.Done()
			<-start
			for index := 0; index < perWorker; index++ {
				todo, err := repository.Create(Todo{Title: "task", Status: StatusOpen})
				if err != nil {
					t.Errorf("Create() error = %v", err)
					return
				}
				createdIDs <- todo.ID
				if _, err := repository.UpdateStatus(todo.ID, StatusCompleted); err != nil {
					t.Errorf("UpdateStatus() error = %v", err)
					return
				}
				if index%2 == 0 {
					if err := repository.Delete(todo.ID); err != nil {
						t.Errorf("Delete() error = %v", err)
						return
					}
				}
			}
		}(worker)
	}

	listDone.Add(1)
	go func() {
		defer listDone.Done()
		<-start
		for {
			select {
			case <-finished:
				return
			default:
				if _, err := repository.List(FilterAll); err != nil {
					t.Errorf("List() error = %v", err)
					return
				}
			}
		}
	}()

	close(start)
	workersDone.Wait()
	close(finished)
	listDone.Wait()
	close(createdIDs)

	seen := make(map[int64]bool, workers*perWorker)
	for id := range createdIDs {
		if id < 1 || seen[id] {
			t.Errorf("created ID %d is invalid or duplicated", id)
		}
		seen[id] = true
	}
	if len(seen) != workers*perWorker {
		t.Errorf("created %d unique IDs, want %d", len(seen), workers*perWorker)
	}

	all, err := repository.List(FilterAll)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != workers*perWorker/2 {
		t.Errorf("List(all) returned %d todos, want %d remaining", len(all), workers*perWorker/2)
	}
}
