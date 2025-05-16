package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	} `json:"user"`
}

type SignInHandler struct {
	userService services.UserService
}

func (h *SignInHandler) Handle(w http.ResponseWriter, r *http.Request, action framework.HasuraAction) {
	var wrapper SignInInputWrapper
	if err := json.Unmarshal(action.Input, &wrapper); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	input := wrapper.Arg1
	if input.Username == "" || input.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	token, user, err := h.userService.SignIn(input.Username, input.Password)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response := SignInResponse{
		Token: token,
	}
	response.User.ID = fmt.Sprint(user.ID)
	response.User.Username = user.Username
	utils.EncodeJSON(w, response)
}

func RegisterSignInHandler(userService services.UserService) {
	dispatcher := framework.GetActionDispatcher(&DefaultHandler{})
	dispatcher.RegisterHandler("signin", &SignInHandler{userService: userService})
}
