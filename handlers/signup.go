package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"app/framework"
	"app/services"
	"app/utils"
)

type SignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpInputWrapper struct {
	Arg1 SignUpInput `json:"arg1"`
}

type SignUpResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type SignUpHandler struct {
	userService services.UserService
}

func (h *SignUpHandler) Handle(w http.ResponseWriter, r *http.Request, action framework.HasuraAction) {
	var wrapper SignUpInputWrapper
	if err := json.Unmarshal(action.Input, &wrapper); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	input := wrapper.Arg1
	if input.Username == "" || input.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	user, err := h.userService.SignUp(input.Username, input.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := SignUpResponse{
		ID:       fmt.Sprint(user.ID),
		Username: user.Username,
	}
	utils.EncodeJSON(w, response)
}

func RegisterSignUpHandler(userService services.UserService) {
	dispatcher := framework.GetActionDispatcher(&DefaultHandler{})
	dispatcher.RegisterHandler("signup", &SignUpHandler{userService: userService})
}
