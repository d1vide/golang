package app

import (
	"net/http"

	"github.com/d1vide/myapp/internal/app/handlers"
	"github.com/d1vide/myapp/internal/app/middlewares"
	"github.com/d1vide/myapp/utils"
)

type pingResp struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Root)
	mux.HandleFunc("/ping", handlers.Ping)
	mux.HandleFunc("/fail", handlers.Fail)

	handler := middlewares.WithRequestID(mux)

	utils.LogInfo("Server is starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		utils.LogError("server error: " + err.Error())
	}
}
