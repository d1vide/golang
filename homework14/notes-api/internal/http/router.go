package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"example.com/notes-api/internal/http/handlers"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/notes", func(r chi.Router) {
			r.Get("/", h.GetNotesV1)

			r.Get("/search", h.SearchNotesV1)

			r.Get("/batch", h.GetNotesBatchV1)

			r.Post("/", h.CreateNote)
			r.Get("/{id}", h.GetNoteByIDV1)
			r.Put("/{id}", h.UpdateNoteV1)
			r.Delete("/{id}", h.DeleteNoteV1)
		})
	})

	return r
}
