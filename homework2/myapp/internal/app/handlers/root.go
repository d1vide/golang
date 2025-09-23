package handlers

import (
	"fmt"
	"net/http"

	"github.com/d1vide/myapp/utils"
)

func Root(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Hello, Go project structure!")
}
