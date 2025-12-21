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
	"example.com/notes-api/internal/db"
	httpLocal "example.com/notes-api/internal/http"
	"example.com/notes-api/internal/http/handlers"
	"example.com/notes-api/internal/repo"
)

func main() {
	connStr := "postgresql://user:password@localhost:5433/notes?sslmode=disable"

	pool, err := db.NewDBPool(connStr)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	log.Println("Successfully connected to database")

	noteRepo := repo.NewNoteDBRepo(pool)

	noteService := service.NewNoteService(noteRepo)

	h := &handlers.Handler{Service: noteService}

	r := httpLocal.NewRouter(h)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
