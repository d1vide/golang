package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"example.com/pz5-security/internal/config"
	"example.com/pz5-security/internal/httpapi"
	"example.com/pz5-security/internal/student"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := student.NewRepo(db)

	stmt, err := repo.PrepareGetByID()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	handler := httpapi.NewHandler(repo, stmt)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/students", handler.GetStudentByID)
	mux.HandleFunc("/students/by-email", handler.GetStudentByEmail)

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      cfg.Addr,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.Println("HTTPS server started on https://localhost:8443")

	go func() {
		log.Println("HTTP redirect server started on http://localhost:8080")

		err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://localhost:8443" + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}))
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
