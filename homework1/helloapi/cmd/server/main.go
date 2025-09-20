package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(
		user{
			ID:   uuid.NewString(),
			Name: "Gopher",
		},
	)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(
		map[string]string{
			"status": "ok",
			"time":   string(time.Now().Format(time.RFC3339)),
		},
	)
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
		log.Printf("APP_PORT is not set, using default port %s", port)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("GET /user", userHandler)
	mux.HandleFunc("GET /health", healthHandler)

	addr := ":" + port
	log.Printf("Starting on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
