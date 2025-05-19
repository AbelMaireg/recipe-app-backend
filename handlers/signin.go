package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"app/framework"
	"app/services"
	"app/utils"
)

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInInputWrapper struct {
	Arg1 SignInInput `json:"arg1"`
}

type SignInResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Bio      string `json:"bio"`
	} `json:"user"`
}

type SignInHandler struct {
	userService services.UserService
}

func (h *SignInHandler) Handle(w http.ResponseWriter, r *http.Request, action framework.HasuraAction) {
	var wrapper SignInInputWrapper
	if err := json.Unmarshal(action.Input, &wrapper); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid input format: "+err.Error())
		return
	}

	input := wrapper.Arg1
	if input.Username == "" || input.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "MISSING_REQUIRED_FIELDS", "Username and password are required")
		return
	}

	token, user, err := h.userService.SignIn(input.Username, input.Password)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no user") {
			utils.WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password")
		} else if strings.Contains(err.Error(), "password") {
			utils.WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to sign in: "+err.Error())
		}
		return
	}

	response := SignInResponse{Token: token}
	response.User.ID = user.ID.String()
	response.User.Username = user.Username
	response.User.Name = user.Name
	response.User.Bio = user.Bio
	utils.EncodeJSON(w, response)
}

func RegisterSignInHandler(userService services.UserService) {
	dispatcher := framework.GetActionDispatcher(&DefaultHandler{})
	dispatcher.RegisterHandler("signin", &SignInHandler{userService: userService})
}
