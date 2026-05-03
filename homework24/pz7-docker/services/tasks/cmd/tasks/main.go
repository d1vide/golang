package main

import (
	"log"
	"net/http"
	"os"
	"example.com/tasks/internal/api"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8083"
	}

	mux := api.GetRouter()

	addr := ":" + port
	log.Println("tasks service started on", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
