package handlers

import (
	"log"

	"app/config"
	"app/framework"
	"app/repositories"
	"app/services"
)

func SetupRoutes(router *framework.Router) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	RegisterSignUpHandler(userService)
	RegisterSignInHandler(userService)

	router.AddPostHandler("/actions", framework.GetActionDispatcher(&DefaultHandler{}).Handle)
	router.AddPostHandler("/events", HandleEvents)
}
