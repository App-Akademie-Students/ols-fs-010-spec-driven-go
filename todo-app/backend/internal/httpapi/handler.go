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
	Create(title string) (todo.Todo, error)
	List(filter todo.Filter) ([]todo.Todo, error)
	UpdateStatus(id int64, status todo.Status) (todo.Todo, error)
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
		if len(values) != 1 {
			writeError(w, http.StatusBadRequest, "Ungültiger Statusfilter.")
			return
		}
		filter = todo.Filter(values[0])
	}

	todos, err := api.service.List(filter)
	if err != nil {
		if errors.Is(err, todo.ErrInvalidFilter) {
			writeError(w, http.StatusBadRequest, "Ungültiger Statusfilter.")
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
		Title string `json:"title"`
	}](w, r)
	if err != nil {
		writeDecodeError(w, err)
		return
	}

	created, err := api.service.Create(request.Title)
	if err != nil {
		if errors.Is(err, todo.ErrInvalidTitle) {
			writeError(w, http.StatusBadRequest, "Titel ist erforderlich. Keine Aufgabe erstellt.")
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
		Status todo.Status `json:"status"`
	}](w, r)
	if err != nil {
		writeDecodeError(w, err)
		return
	}

	updated, err := api.service.UpdateStatus(id, request.Status)
	if err != nil {
		switch {
		case errors.Is(err, todo.ErrInvalidStatus):
			writeError(w, http.StatusBadRequest, "Ungültiger Status.")
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
