package handlers

import (
	"net/http"

	"go-graphql-app/framework"
	"go-graphql-app/utils"
)

type DefaultHandler struct{}

func (h *DefaultHandler) Handle(w http.ResponseWriter, r *http.Request, action framework.HasuraAction) {
	utils.WriteError(w, http.StatusBadRequest, "Unknown action")
}
