package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"example.com/notes-api/internal/core/service"
	"example.com/notes-api/internal/http/dto"
	"example.com/notes-api/internal/repo"
)

type Handler struct {
	Service *service.NoteService
}

// CreateNote godoc
// @Summary      Создать заметку
// @Description  Создаёт новую заметку и возвращает объект
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        note  body      dto.CreateNoteRequest  true  "Данные заметки"
// @Success      201   {object}  dto.NoteResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes [post]
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

// GetNotes godoc
// @Summary      Список заметок
// @Description  Возвращает список всех заметок
// @Tags         notes
// @Produce      json
// @Success      200  {array}   core.Note
// @Failure      500  {object}  map[string]string
// @Router       /notes [get]
func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}

// GetNoteByID godoc
// @Summary      Получить заметку
// @Description  Возвращает заметку по идентификатору
// @Tags         notes
// @Produce      json
// @Param        id   path      int64  true  "ID заметки"
// @Success      200  {object}  core.Note
// @Failure      400  {object}  map[string]string  "Некорректный ID"
// @Failure      404  {object}  map[string]string  "Заметка не найдена"
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [get]
func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

// UpdateNote godoc
// @Summary      Обновить заметку
// @Description  Полностью обновляет заметку по ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id    path      int64                 true  "ID заметки"
// @Param        note  body      dto.UpdateNoteRequest true  "Новые данные заметки"
// @Success      200   {object}  dto.NoteResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes/{id} [put]
func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
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

// DeleteNote godoc
// @Summary      Удалить заметку
// @Description  Удаляет заметку по ID
// @Tags         notes
// @Param        id   path  int64  true  "ID заметки"
// @Success      204  "Заметка удалена"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
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
