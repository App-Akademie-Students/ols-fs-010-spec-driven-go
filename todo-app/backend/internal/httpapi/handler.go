package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/App-Akademie-Students/ols-fs-010-spec-driven-go/todo-app/backend/internal/todo"
)

const maxRequestBodyBytes = 1 << 20

type todoService interface {
	Create(title string, priority ...todo.Priority) (todo.Todo, error)
	ListWithFilters(filter todo.Filter, priority todo.PriorityFilter) ([]todo.Todo, error)
	Update(id int64, status *todo.Status, priority *todo.Priority) (todo.Todo, error)
	Delete(id int64) error
}

type handler struct {
	service todoService
}

func NewHandler(service todoService) http.Handler {
	api := &handler{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/todos", api.list)
	mux.HandleFunc("POST /api/todos", api.create)
	mux.HandleFunc("PATCH /api/todos/{id}", api.updateStatus)
	mux.HandleFunc("DELETE /api/todos/{id}", api.delete)
	return mux
}

func (api *handler) list(w http.ResponseWriter, r *http.Request) {
	filter := todo.FilterAll
	if values, exists := r.URL.Query()["status"]; exists {
		if len(values) != 1 || values[0] == "" {
			writeError(w, http.StatusBadRequest, "Ungültiger Statusfilter.")
			return
		}
		filter = todo.Filter(values[0])
	}
	priorityFilter := todo.PriorityFilterAll
	if values, exists := r.URL.Query()["priority"]; exists {
		if len(values) != 1 || values[0] == "" {
			writeError(w, http.StatusBadRequest, "Ungültiger Prioritätsfilter.")
			return
		}
		priorityFilter = todo.PriorityFilter(values[0])
	}

	todos, err := api.service.ListWithFilters(filter, priorityFilter)
	if err != nil {
		if errors.Is(err, todo.ErrInvalidFilter) {
			writeError(w, http.StatusBadRequest, "Ungültiger Statusfilter.")
			return
		}
		if errors.Is(err, todo.ErrInvalidPriorityFilter) {
			writeError(w, http.StatusBadRequest, "Ungültiger Prioritätsfilter.")
			return
		}
		writeError(w, http.StatusInternalServerError, "Aktion fehlgeschlagen. Bitte erneut versuchen.")
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Todos []todo.Todo `json:"todos"`
	}{Todos: todos})
}

func (api *handler) create(w http.ResponseWriter, r *http.Request) {
	request, err := decodeJSON[struct {
		Title    string          `json:"title"`
		Priority json.RawMessage `json:"priority"`
	}](w, r)
	if err != nil {
		writeDecodeError(w, err)
		return
	}

	var created todo.Todo
	if len(request.Priority) == 0 {
		created, err = api.service.Create(request.Title)
	} else {
		var priority todo.Priority
		if string(request.Priority) == "null" || json.Unmarshal(request.Priority, &priority) != nil {
			writeError(w, http.StatusBadRequest, "Ungültige Priorität.")
			return
		}
		created, err = api.service.Create(request.Title, priority)
	}
	if err != nil {
		if errors.Is(err, todo.ErrInvalidTitle) {
			writeError(w, http.StatusBadRequest, "Titel ist erforderlich. Keine Aufgabe erstellt.")
			return
		}
		if errors.Is(err, todo.ErrInvalidPriority) {
			writeError(w, http.StatusBadRequest, "Ungültige Priorität.")
			return
		}
		writeError(w, http.StatusInternalServerError, "Aktion fehlgeschlagen. Bitte erneut versuchen.")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (api *handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Ungültige Aufgaben-ID.")
		return
	}
	request, err := decodeJSON[struct {
		Status   json.RawMessage `json:"status"`
		Priority json.RawMessage `json:"priority"`
	}](w, r)
	if err != nil {
		writeDecodeError(w, err)
		return
	}

	if len(request.Status) == 0 && len(request.Priority) == 0 {
		writeError(w, http.StatusBadRequest, "Mindestens ein Aufgabenfeld muss angegeben werden.")
		return
	}
	var status *todo.Status
	if len(request.Status) != 0 {
		var value todo.Status
		if string(request.Status) == "null" || json.Unmarshal(request.Status, &value) != nil {
			writeError(w, http.StatusBadRequest, "Ungültiger Status.")
			return
		}
		status = &value
	}
	var priority *todo.Priority
	if len(request.Priority) != 0 {
		var value todo.Priority
		if string(request.Priority) == "null" || json.Unmarshal(request.Priority, &value) != nil {
			writeError(w, http.StatusBadRequest, "Ungültige Priorität.")
			return
		}
		priority = &value
	}

	updated, err := api.service.Update(id, status, priority)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrInvalidStatus):
			writeError(w, http.StatusBadRequest, "Ungültiger Status.")
		case errors.Is(err, todo.ErrInvalidPriority):
			writeError(w, http.StatusBadRequest, "Ungültige Priorität.")
		case errors.Is(err, todo.ErrEmptyUpdate):
			writeError(w, http.StatusBadRequest, "Mindestens ein Aufgabenfeld muss angegeben werden.")
		case errors.Is(err, todo.ErrNotFound):
			writeError(w, http.StatusNotFound, "Aufgabe ist nicht mehr verfügbar.")
		default:
			writeError(w, http.StatusInternalServerError, "Aktion fehlgeschlagen. Bitte erneut versuchen.")
		}
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (api *handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Ungültige Aufgaben-ID.")
		return
	}
	if err := api.service.Delete(id); err != nil {
		if errors.Is(err, todo.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Aufgabe ist nicht mehr verfügbar.")
			return
		}
		writeError(w, http.StatusInternalServerError, "Aktion fehlgeschlagen. Bitte erneut versuchen.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var request T
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return request, errors.New("request must contain exactly one JSON value")
		}
		return request, err
	}
	return request, nil
}

func writeDecodeError(w http.ResponseWriter, err error) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		writeError(w, http.StatusRequestEntityTooLarge, "Anfrage ist zu groß.")
		return
	}
	writeError(w, http.StatusBadRequest, "Ungültige Anfrage.")
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid todo ID %q", value)
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode HTTP response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, struct {
		Error string `json:"error"`
	}{Error: message})
}
