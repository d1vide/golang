package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"example.com/notes-api/internal/core"
	"example.com/notes-api/internal/core/service"
	"example.com/notes-api/internal/http/dto"
	"example.com/notes-api/internal/repo"
)

type Handler struct {
	Service *service.NoteService
}

type PaginatedResponse struct {
	Notes  []core.Note `json:"notes"`
	Total  int64       `json:"total,omitempty"`
	Cursor struct {
		ID        int64     `json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
	} `json:"cursor,omitempty"`
	HasMore bool `json:"has_more"`
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	note, err := h.Service.Create(req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetNotesV1 - обновленный метод с поддержкой разных видов пагинации
// @Summary      Получить заметки с пагинацией
// @Description  Возвращает заметки с offset или keyset пагинацией
// @Tags         notes
// @Produce      json
// @Param        limit   query   int     false  "Лимит записей (по умолчанию 20)"
// @Param        offset  query   int64   false  "Смещение (для offset пагинации)"
// @Param        mode    query   string  false  "Режим пагинации: offset или keyset (по умолчанию offset)"
// @Param        cursor_id       query   int64   false  "ID курсора (для keyset пагинации)"
// @Param        cursor_time     query   string  false  "Время курсора в RFC3339 (для keyset пагинации)"
// @Success      200  {object}  PaginatedResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/notes [get]
func (h *Handler) GetNotesV1(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	offset, _ := strconv.ParseInt(query.Get("offset"), 10, 64)

	mode := query.Get("mode")
	if mode == "" {
		mode = "offset"
	}

	var cursorID int64
	var cursorTime time.Time

	if mode == "keyset" {
		cursorID, _ = strconv.ParseInt(query.Get("cursor_id"), 10, 64)

		if cursorTimeStr := query.Get("cursor_time"); cursorTimeStr != "" {
			if parsedTime, err := time.Parse(time.RFC3339, cursorTimeStr); err == nil {
				cursorTime = parsedTime
			}
		}
	}

	notes, err := h.Service.GetNotesPaginated(limit, offset, cursorID, cursorTime, mode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := PaginatedResponse{
		Notes:   notes,
		HasMore: len(notes) >= limit,
	}

	if mode == "keyset" && len(notes) > 0 {
		lastNote := notes[len(notes)-1]
		resp.Cursor.ID = lastNote.ID
		resp.Cursor.CreatedAt = lastNote.CreatedAt
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GetNotesBatchV1 - получение заметок батчем
// @Summary      Получить заметки батчем
// @Description  Возвращает несколько заметок по их ID
// @Tags         notes
// @Produce      json
// @Param        ids   query   string  true  "Список ID через запятую (например: 1,2,3,4,5)"
// @Success      200  {object}  map[string]core.Note
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/notes/batch [get]
func (h *Handler) GetNotesBatchV1(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	idsParam := query.Get("ids")
	if idsParam == "" {
		http.Error(w, "ids parameter is required", http.StatusBadRequest)
		return
	}

	idStrs := strings.Split(idsParam, ",")
	ids := make([]int64, 0, len(idStrs))

	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid ID: %s", idStr), http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	notes, err := h.Service.GetNotesBatch(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}

// GetNoteByIDV1 - получение одной заметки (версия v1)
// @Summary      Получить заметку по ID
// @Description  Возвращает заметку по её ID
// @Tags         notes
// @Produce      json
// @Param        id   path      int64  true  "ID заметки"
// @Success      200  {object}  core.Note
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/notes/{id} [get]
func (h *Handler) GetNoteByIDV1(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	note, err := h.Service.GetByID(id)
	if err == repo.ErrNotFound {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

// SearchNotesV1 - поиск заметок
// @Summary      Поиск заметок по заголовку
// @Description  Возвращает заметки, соответствующие поисковому запросу
// @Tags         notes
// @Produce      json
// @Param        q      query   string  true  "Поисковый запрос"
// @Param        limit  query   int     false  "Лимит результатов (по умолчанию 20)"
// @Success      200  {array}   core.Note
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/notes/search [get]
func (h *Handler) SearchNotesV1(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	searchQuery := query.Get("q")
	if searchQuery == "" {
		http.Error(w, "q parameter is required", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	notes, err := h.Service.SearchNotes(searchQuery, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}

func (h *Handler) UpdateNoteV1(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	note, err := h.Service.Update(id, req.Title, req.Content)
	if err == repo.ErrNotFound {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) DeleteNoteV1(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.Delete(id); err == repo.ErrNotFound {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
