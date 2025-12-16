package httpx

import (
	"net/http"

	"example.com/notes-api/docs"
	"example.com/notes-api/internal/http/handlers"
	"example.com/notes-api/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Post("/", h.CreateNote)
		r.Get("/", h.GetNotes)
		r.Get("/{id}", h.GetNoteByID)
		r.Put("/{id}", h.UpdateNote)
		r.Delete("/{id}", h.DeleteNote)
	})

	r.Get("/docs/*", httpSwagger.WrapHandler)

	r.Get("/redoc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "internal/http/static/redoc.html")
	})
	r.Get("/docs/swagger.json", serveSwaggerJSON)

	return r
}

func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b := docs.SwaggerInfo.ReadDoc()
	w.Write([]byte(b))
}
