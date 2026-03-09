package taskshttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	"tech-ip-sem2/services/tasks/internal/service"
)

type Handler struct {
	svc  *service.TasksService
	auth *authclient.Client
}

func NewHandler(svc *service.TasksService, auth *authclient.Client) *Handler {
	return &Handler{svc: svc, auth: auth}
}

func (h *Handler) authGuard(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSON(w, http.StatusUnauthorized, errBody("missing authorization header"))
		return false
	}

	_, err := h.auth.Verify(r.Context(), authHeader)
	if err == nil {
		return true
	}

	if errors.Is(err, authclient.ErrUnauthorized) {
		writeJSON(w, http.StatusUnauthorized, errBody("unauthorized"))
		return false
	}

	writeJSON(w, http.StatusServiceUnavailable, errBody("auth service unavailable"))
	return false
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	if !h.authGuard(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errBody("method not allowed"))
	}
}

func (h *Handler) Task(w http.ResponseWriter, r *http.Request) {
	if !h.authGuard(w, r) {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errBody("missing task id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPatch:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errBody("method not allowed"))
	}
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.svc.List()
	if tasks == nil {
		tasks = []service.TaskSummary{}
	}
	writeJSON(w, http.StatusOK, tasks)
}

type createRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Title) == "" {
		writeJSON(w, http.StatusBadRequest, errBody("bad request: title is required"))
		return
	}

	task := h.svc.Create(service.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	})
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.svc.Get(id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, errBody("task not found"))
		return
	}
	writeJSON(w, http.StatusOK, task)
}

type updateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
	Done        *bool   `json:"done"`
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errBody("bad request"))
		return
	}

	task, err := h.svc.Update(id, service.UpdateInput{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Done:        req.Done,
	})
	if errors.Is(err, service.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, errBody("task not found"))
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	err := h.svc.Delete(id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, errBody("task not found"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func errBody(msg string) map[string]string {
	return map[string]string{"error": msg}
}
