package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pz16/internal/db"
	"pz16/internal/httpapi"
	"pz16/internal/repo"
	"pz16/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	serverPort := "8080"
	connStr := "host=localhost port=54321 user=test password=test dbname=notes_test sslmode=disable"

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer sqlDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	db.MustApplyMigrations(sqlDB)
	log.Println("Migrations applied successfully")

	noteRepo := repo.NoteRepo{DB: sqlDB}
	svc := service.Service{Notes: noteRepo}
	router := httpapi.Router{Svc: &svc}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	router.Register(r)

	r.GET("/health", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", serverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
