package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"app/framework"
	"app/services"
	"app/utils"
)

type SignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
}

type SignUpInputWrapper struct {
	Arg1 SignUpInput `json:"arg1"`
}

type SignUpResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
}

type SignUpHandler struct {
	userService services.UserService
}

func (h *SignUpHandler) Handle(w http.ResponseWriter, r *http.Request, action framework.HasuraAction) {
	var wrapper SignUpInputWrapper
	if err := json.Unmarshal(action.Input, &wrapper); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid input format: "+err.Error())
		return
	}

	input := wrapper.Arg1
	if input.Username == "" || input.Password == "" || input.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "MISSING_REQUIRED_FIELDS", "Username, password, and name are required")
		return
	}

	// Basic password validation
	if len(input.Password) < 8 {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PASSWORD", "Password must be at least 8 characters long")
		return
	}

	user, err := h.userService.SignUp(input.Username, input.Password, input.Name, input.Bio)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "username") {
			utils.WriteError(w, http.StatusBadRequest, "USERNAME_TAKEN", "Username is already taken")
		} else if strings.Contains(err.Error(), "password") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_PASSWORD", "Invalid password format")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to sign up: "+err.Error())
		}
		return
	}

	response := SignUpResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Bio:      user.Bio,
	}
	utils.EncodeJSON(w, response)
}

func RegisterSignUpHandler(userService services.UserService) {
	dispatcher := framework.GetActionDispatcher(&DefaultHandler{})
	dispatcher.RegisterHandler("signup", &SignUpHandler{userService: userService})
}
