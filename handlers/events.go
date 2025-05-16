package handlers

import (
	"log"
	"net/http"

	"go-graphql-app/models"
	"go-graphql-app/utils"
)

type UserEventPayload struct {
	Event struct {
		Op   string `json:"op"`
		Data struct {
			New models.User `json:"new"`
		} `json:"data"`
	} `json:"event"`
}

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	var payload UserEventPayload
	if err := utils.DecodeJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid event payload")
		return
	}

	if payload.Event.Op == "INSERT" {
		log.Printf("User created: ID=%d, Username=%s", payload.Event.Data.New.ID, payload.Event.Data.New.Username)
	}

	w.WriteHeader(http.StatusOK)
}
