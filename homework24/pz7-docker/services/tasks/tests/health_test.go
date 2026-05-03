package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"example.com/tasks/internal/api"
)

func TestHealth(t *testing.T) {
	router := api.GetRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Код ответа не 200, а %d", rr.Code)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("Ожидался статус ok, получили %s", response["status"])
	}
}
