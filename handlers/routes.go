package handlers

import (
	"app/framework"
	"app/services"
)

func SetupRoutes(router *framework.Router, userService services.UserService) {
	RegisterSignUpHandler(userService)
	RegisterSignInHandler(userService)

	router.AddPostHandler("/actions", framework.GetActionDispatcher(&DefaultHandler{}).Handle)
	router.AddPostHandler("/events", HandleEvents)
}
