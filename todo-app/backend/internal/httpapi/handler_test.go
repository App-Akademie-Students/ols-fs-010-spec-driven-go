package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/App-Akademie-Students/ols-fs-010-spec-driven-go/todo-app/backend/internal/todo"
)

func TestCreateAndListTodos(t *testing.T) {
	handler := NewHandler(todo.NewService(todo.NewInMemoryRepository()))

	createdResponse := perform(handler, http.MethodPost, "/api/todos", `{"title":"  Prepare demo  "}`)
	if createdResponse.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d: %s", createdResponse.Code, http.StatusCreated, createdResponse.Body)
	}
	created := decodeTodo(t, createdResponse)
	if created.Title != "Prepare demo" || created.Status != todo.StatusOpen || created.Priority != todo.PriorityMedium || created.ID < 1 {
		t.Fatalf("created todo = %#v, want trimmed title, positive ID, and open status", created)
	}
	if got := createdResponse.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	secondResponse := perform(handler, http.MethodPost, "/api/todos", `{"title":"Second"}`)
	second := decodeTodo(t, secondResponse)
	allResponse := perform(handler, http.MethodGet, "/api/todos", "")
	if allResponse.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d", allResponse.Code, http.StatusOK)
	}
	var result struct {
		Todos []todo.Todo `json:"todos"`
	}
	if err := json.Unmarshal(allResponse.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(result.Todos) != 2 || result.Todos[0].ID != created.ID || result.Todos[1].ID != second.ID {
		t.Fatalf("GET todos = %#v, want both todos in creation order", result.Todos)
	}

	duplicateResponse := perform(handler, http.MethodPost, "/api/todos", `{"title":"Prepare demo"}`)
	duplicate := decodeTodo(t, duplicateResponse)
	if duplicate.ID == created.ID || duplicate.Title != created.Title {
		t.Fatalf("duplicate todo = %#v, want same title and distinct ID", duplicate)
	}
}

func TestPriorityCreateAndUpdateRequests(t *testing.T) {
	handler := NewHandler(todo.NewService(todo.NewInMemoryRepository()))

	createdResponse := perform(handler, http.MethodPost, "/api/todos", `{"title":"priority task","priority":"low"}`)
	if createdResponse.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d: %s", createdResponse.Code, http.StatusCreated, createdResponse.Body)
	}
	created := decodeTodo(t, createdResponse)
	if created.Priority != todo.PriorityLow {
		t.Fatalf("created priority = %q, want %q", created.Priority, todo.PriorityLow)
	}

	priorityUpdate := perform(handler, http.MethodPatch, "/api/todos/"+itoa(created.ID), `{"priority":"high"}`)
	if priorityUpdate.Code != http.StatusOK {
		t.Fatalf("priority PATCH status = %d, want %d: %s", priorityUpdate.Code, http.StatusOK, priorityUpdate.Body)
	}
	updated := decodeTodo(t, priorityUpdate)
	if updated.Priority != todo.PriorityHigh || updated.Title != created.Title || updated.Status != created.Status {
		t.Fatalf("priority PATCH result = %#v; want only priority changed", updated)
	}

	statusOnlyUpdate := perform(handler, http.MethodPatch, "/api/todos/"+itoa(created.ID), `{"status":"completed"}`)
	if statusOnlyUpdate.Code != http.StatusOK {
		t.Fatalf("status-only PATCH status = %d, want %d", statusOnlyUpdate.Code, http.StatusOK)
	}
	statusOnly := decodeTodo(t, statusOnlyUpdate)
	if statusOnly.Priority != todo.PriorityHigh || statusOnly.Status != todo.StatusCompleted {
		t.Fatalf("status-only PATCH result = %#v; want existing priority retained", statusOnly)
	}

	invalidPriority := perform(handler, http.MethodPatch, "/api/todos/"+itoa(created.ID), `{"priority":"urgent"}`)
	if invalidPriority.Code != http.StatusBadRequest {
		t.Fatalf("invalid priority PATCH status = %d, want %d", invalidPriority.Code, http.StatusBadRequest)
	}
	assertErrorResponse(t, invalidPriority)

	emptyPriority := perform(handler, http.MethodPost, "/api/todos", `{"title":"empty priority","priority":""}`)
	if emptyPriority.Code != http.StatusBadRequest {
		t.Fatalf("empty priority POST status = %d, want %d", emptyPriority.Code, http.StatusBadRequest)
	}
	assertErrorResponse(t, emptyPriority)
}

func TestPriorityFiltersCombineWithStatus(t *testing.T) {
	handler := NewHandler(todo.NewService(todo.NewInMemoryRepository()))
	create := func(title, priority string) todo.Todo {
		t.Helper()
		body := `{"title":"` + title + `","priority":"` + priority + `"}`
		response := perform(handler, http.MethodPost, "/api/todos", body)
		if response.Code != http.StatusCreated {
			t.Fatalf("POST status = %d, want %d: %s", response.Code, http.StatusCreated, response.Body)
		}
		return decodeTodo(t, response)
	}

	lowOpen := create("low open", "low")
	highOpen := create("high open", "high")
	mediumOpen := create("medium open", "medium")
	highCompleted := create("high completed", "high")
	completedResponse := perform(handler, http.MethodPatch, "/api/todos/"+itoa(highCompleted.ID), `{"status":"completed"}`)
	if completedResponse.Code != http.StatusOK {
		t.Fatalf("PATCH status = %d, want %d", completedResponse.Code, http.StatusOK)
	}

	tests := []struct {
		name       string
		path       string
		wantIDs    []int64
		wantStatus int
	}{
		{name: "combined open high", path: "/api/todos?status=open&priority=high", wantIDs: []int64{highOpen.ID}, wantStatus: http.StatusOK},
		{name: "combined completed high", path: "/api/todos?status=completed&priority=high", wantIDs: []int64{highCompleted.ID}, wantStatus: http.StatusOK},
		{name: "priority only", path: "/api/todos?priority=medium", wantIDs: []int64{mediumOpen.ID}, wantStatus: http.StatusOK},
		{name: "empty match", path: "/api/todos?status=completed&priority=low", wantIDs: []int64{}, wantStatus: http.StatusOK},
		{name: "invalid priority", path: "/api/todos?priority=urgent", wantStatus: http.StatusBadRequest},
		{name: "empty priority", path: "/api/todos?priority=", wantStatus: http.StatusBadRequest},
		{name: "empty status", path: "/api/todos?status=", wantStatus: http.StatusBadRequest},
		{name: "repeated priority", path: "/api/todos?priority=high&priority=low", wantStatus: http.StatusBadRequest},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response := perform(handler, http.MethodGet, testCase.path, "")
			if response.Code != testCase.wantStatus {
				t.Fatalf("GET status = %d, want %d: %s", response.Code, testCase.wantStatus, response.Body)
			}
			if testCase.wantStatus != http.StatusOK {
				assertErrorResponse(t, response)
				return
			}
			var result struct {
				Todos []todo.Todo `json:"todos"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if len(result.Todos) != len(testCase.wantIDs) {
				t.Fatalf("GET todos = %#v, want IDs %v", result.Todos, testCase.wantIDs)
			}
			for index, wantID := range testCase.wantIDs {
				if result.Todos[index].ID != wantID {
					t.Errorf("GET todos[%d].ID = %d, want %d", index, result.Todos[index].ID, wantID)
				}
			}
		})
	}
	if lowOpen.ID == 0 {
		t.Fatal("low priority todo was not created")
	}
}

func TestRejectsInvalidCreateRequests(t *testing.T) {
	handler := NewHandler(todo.NewService(todo.NewInMemoryRepository()))
	cases := []struct {
		name string
		body string
	}{
		{name: "blank", body: `{"title":" \u00a0 "}`},
		{name: "malformed JSON", body: `{"title":`},
		{name: "unknown field", body: `{"title":"valid","extra":true}`},
		{name: "trailing JSON", body: `{"title":"valid"} {}`},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			response := perform(handler, http.MethodPost, "/api/todos", testCase.body)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusBadRequest, response.Body)
			}
			assertErrorResponse(t, response)
		})
	}

	listResponse := perform(handler, http.MethodGet, "/api/todos", "")
	var result struct {
		Todos []todo.Todo `json:"todos"`
	}
	if err := json.Unmarshal(listResponse.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Todos) != 0 {
		t.Fatalf("invalid create requests created %d todos", len(result.Todos))
	}
}

func TestRejectsOversizedRequestBody(t *testing.T) {
	handler := NewHandler(todo.NewService(todo.NewInMemoryRepository()))
	body := `{"title":"` + strings.Repeat("x", maxRequestBodyBytes) + `"}`
	response := perform(handler, http.MethodPost, "/api/todos", body)
	if response.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusRequestEntityTooLarge)
	}
	assertErrorResponse(t, response)
}

func TestStatusFiltersAndMutations(t *testing.T) {
	handler := NewHandler(todo.NewService(todo.NewInMemoryRepository()))
	first := decodeTodo(t, perform(handler, http.MethodPost, "/api/todos", `{"title":"first"}`))
	second := decodeTodo(t, perform(handler, http.MethodPost, "/api/todos", `{"title":"second"}`))

	updateResponse := perform(handler, http.MethodPatch, "/api/todos/"+itoa(first.ID), `{"status":"completed"}`)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("PATCH status = %d, want %d: %s", updateResponse.Code, http.StatusOK, updateResponse.Body)
	}
	updated := decodeTodo(t, updateResponse)
	if updated.Status != todo.StatusCompleted {
		t.Fatalf("updated status = %q, want %q", updated.Status, todo.StatusCompleted)
	}

	for _, testCase := range []struct {
		filter string
		wantID int64
	}{
		{filter: "completed", wantID: first.ID},
		{filter: "open", wantID: second.ID},
	} {
		response := perform(handler, http.MethodGet, "/api/todos?status="+testCase.filter, "")
		if response.Code != http.StatusOK {
			t.Fatalf("GET filter %q status = %d", testCase.filter, response.Code)
		}
		var result struct {
			Todos []todo.Todo `json:"todos"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if len(result.Todos) != 1 || result.Todos[0].ID != testCase.wantID {
			t.Errorf("filter %q returned %#v, want ID %d", testCase.filter, result.Todos, testCase.wantID)
		}
	}

	reopenResponse := perform(handler, http.MethodPatch, "/api/todos/"+itoa(first.ID), `{"status":"open"}`)
	if got := decodeTodo(t, reopenResponse).Status; got != todo.StatusOpen {
		t.Errorf("reopened status = %q, want %q", got, todo.StatusOpen)
	}
	deleteResponse := perform(handler, http.MethodDelete, "/api/todos/"+itoa(second.ID), "")
	if deleteResponse.Code != http.StatusNoContent || deleteResponse.Body.Len() != 0 {
		t.Fatalf("DELETE response = %d %q, want 204 and empty body", deleteResponse.Code, deleteResponse.Body)
	}
	if response := perform(handler, http.MethodGet, "/api/todos", ""); strings.Contains(response.Body.String(), `"id":`+itoa(second.ID)) {
		t.Errorf("deleted todo remains in list response: %s", response.Body)
	}
}

func TestRejectsInvalidFiltersStatusesAndUnknownIDs(t *testing.T) {
	handler := NewHandler(todo.NewService(todo.NewInMemoryRepository()))
	for _, path := range []string{"/api/todos?status=archived", "/api/todos?status=open&status=completed"} {
		response := perform(handler, http.MethodGet, path, "")
		if response.Code != http.StatusBadRequest {
			t.Errorf("GET %s status = %d, want %d", path, response.Code, http.StatusBadRequest)
		}
		assertErrorResponse(t, response)
	}

	created := decodeTodo(t, perform(handler, http.MethodPost, "/api/todos", `{"title":"task"}`))
	response := perform(handler, http.MethodPatch, "/api/todos/"+itoa(created.ID), `{"status":"archived"}`)
	if response.Code != http.StatusBadRequest {
		t.Errorf("invalid status response = %d, want %d", response.Code, http.StatusBadRequest)
	}
	assertErrorResponse(t, response)

	for _, operation := range []struct {
		method string
		body   string
	}{
		{method: http.MethodPatch, body: `{"status":"open"}`},
		{method: http.MethodDelete},
	} {
		response := perform(handler, operation.method, "/api/todos/9999", operation.body)
		if response.Code != http.StatusNotFound {
			t.Errorf("%s missing todo response = %d, want %d", operation.method, response.Code, http.StatusNotFound)
		}
		var payload map[string]string
		if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil || payload["error"] != "Aufgabe ist nicht mehr verfügbar." {
			t.Errorf("missing todo error = %v, %s", err, response.Body)
		}
	}

	badID := perform(handler, http.MethodDelete, "/api/todos/not-a-number", "")
	if badID.Code != http.StatusBadRequest {
		t.Errorf("invalid ID response = %d, want %d", badID.Code, http.StatusBadRequest)
	}
	assertErrorResponse(t, badID)
}

func perform(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func decodeTodo(t *testing.T, response *httptest.ResponseRecorder) todo.Todo {
	t.Helper()
	var result todo.Todo
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode todo response %q: %v", response.Body, err)
	}
	return result
}

func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	var result map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode error response %q: %v", response.Body, err)
	}
	if result["error"] == "" {
		t.Errorf("error response has no message: %s", response.Body)
	}
}

func itoa(value int64) string {
	return strconv.FormatInt(value, 10)
}
