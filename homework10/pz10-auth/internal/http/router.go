package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"example.com/pz10-auth/internal/core"
	"example.com/pz10-auth/internal/http/middleware"
	"example.com/pz10-auth/internal/platform/config"
	"example.com/pz10-auth/internal/platform/jwt"
	"example.com/pz10-auth/internal/repo"
)

func Build(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	// DI
	userRepo := repo.NewUserMem() // храним заранее захэшированных юзеров (email, bcrypt)
	jwtv := jwt.NewRS256()
	svc := core.NewService(userRepo, jwtv)

	// Общие публичные маршруты
	r.Post("/api/v1/login", svc.LoginHandler)
	r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthN(jwtv)(http.HandlerFunc(svc.MeHandler)).ServeHTTP(w, r)
	})

	// Защищённые маршруты для всех пользователей
	r.Group(func(priv chi.Router) {
		priv.Use(middleware.AuthN(jwtv))
		priv.Use(middleware.AuthZRoles("user", "admin"))
		priv.Get("/api/v1/users/{id}", svc.UserByID) // ABAC
	})

	// Только админ
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthN(jwtv))
		admin.Use(middleware.AuthZRoles("admin"))
		admin.Route("/api/v1/admin", func(r chi.Router) {
			r.Get("/stats", svc.AdminStats)
		})
	})

	r.Post("/api/v1/refresh", svc.RefreshHandler)

	return r
}
