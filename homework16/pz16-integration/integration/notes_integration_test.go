package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"pz16/internal/db"
	"pz16/internal/httpapi"
	"pz16/internal/repo"
	"pz16/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func newServer(t *testing.T, dsn string) *httptest.Server {
	t.Helper()
	dbx, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	db.MustApplyMigrations(dbx)

	r := gin.Default()
	svc := service.Service{Notes: repo.NoteRepo{DB: dbx}}
	httpapi.Router{Svc: &svc}.Register(r)

	return httptest.NewServer(r)
}

func TestCreateAndGetNote(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("DB_DSN not set (use `make up` and `make test`)")
	}
	srv := newServer(t, dsn)
	defer srv.Close()

	// 1) Create
	resp, err := http.Post(srv.URL+"/notes", "application/json",
		strings.NewReader(`{"title":"Hello","content":"World"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("status %d != 201, body: %s", resp.StatusCode, string(body))
	}

	var createdNote struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	id := createdNote.ID

	// 2) Get
	resp2, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		t.Fatalf("status %d != 200, body: %s", resp2.StatusCode, string(body))
	}

	// 3) Verify data
	var retrievedNote struct {
		ID      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&retrievedNote); err != nil {
		t.Fatal(err)
	}

	if retrievedNote.ID != id {
		t.Errorf("expected id %d, got %d", id, retrievedNote.ID)
	}
	if retrievedNote.Title != "Hello" {
		t.Errorf("expected title 'Hello', got '%s'", retrievedNote.Title)
	}
	if retrievedNote.Content != "World" {
		t.Errorf("expected content 'World', got '%s'", retrievedNote.Content)
	}
}
