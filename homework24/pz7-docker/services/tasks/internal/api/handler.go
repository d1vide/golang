package api

import (
	"encoding/json"
	"net/http"
)

func GetRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)
	return mux
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "tasks",
	})
}
