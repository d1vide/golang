package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/pz3-http/internal/api"
	"example.com/pz3-http/internal/storage"
)

func RunServerWithGracefulShutdown(handler http.Handler, addr string) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err

	case <-osSignals:
		log.Println("Starting graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			return server.Close()
		}

		log.Println("Server stopped gracefully")
		return nil
	}
}

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Коллекция
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	// Элемент
	mux.HandleFunc("GET /tasks/", h.GetTask)

	mux.HandleFunc("DELETE /tasks/", h.DeleteTask)
	mux.HandleFunc("PATCH /tasks/", h.PatchTask)

	// Подключаем middlewares
	middlewares := api.CORS(api.Logging(mux))

	addr := ":8080"
	log.Println("listening on", addr)
	if err := RunServerWithGracefulShutdown(middlewares, addr); err != nil {
		log.Fatal(err)
	}
}
