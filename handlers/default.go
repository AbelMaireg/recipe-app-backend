package handlers

import (
	"net/http"

	"app/framework"
	"app/utils"
)

type DefaultHandler struct{}

func (h *DefaultHandler) Handle(w http.ResponseWriter, r *http.Request, action framework.HasuraAction) {
	utils.WriteError(w, http.StatusBadRequest, "UNKNOWN_ACTION", "Unknown action")
}
