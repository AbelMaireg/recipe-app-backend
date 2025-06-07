package handlers

import (
	"log"
	"net/http"
)

type HealthCheckHandler struct{}

func (h *HealthCheckHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check request received.")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
