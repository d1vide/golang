// Package main Notes API server.
//
// @title           Notes API
// @version         1.0
// @description     Учебный REST API для заметок (CRUD).
// @contact.name    Backend Course
// @contact.email   example@university.ru
// @BasePath        /api/v1
package main

import (
	"log"
	"net/http"

	_ "example.com/notes-api/docs"

	"example.com/notes-api/internal/core/service"
	httpLocal "example.com/notes-api/internal/http"
	"example.com/notes-api/internal/http/handlers"
	"example.com/notes-api/internal/repo"
)

func main() {
	repo := repo.NewNoteRepoMem()
	service := service.NewNoteService(repo)
	h := &handlers.Handler{Service: service}

	r := httpLocal.NewRouter(h)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
