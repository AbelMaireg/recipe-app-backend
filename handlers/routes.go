package handlers

import (
	"log"

	"app/config"
	"app/framework"
	"app/repositories"
	"app/services"
)

func SetupRoutes(router *framework.Router) {
	cfg, minioCfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	minioClient, err := minioCfg.GetClient()
	if err != nil {
		log.Fatal("Failed to initialize MinIO client:", err)
	}

	userService := services.NewUserService(repositories.NewUserRepository(db))
	recipeService := services.NewRecipeService(repositories.NewRecipeRepository(db))
	recipePictureUploadHandler := NewUploadRecipePictureHandler(recipeService, minioClient, minioCfg.Bucket)

	RegisterSignUpHandler(userService)
	RegisterSignInHandler(userService)

	router.AddPostHandler("/actions", framework.GetActionDispatcher(&DefaultHandler{}).Handle)
	router.AddPostHandler("/events", HandleEvents)
	router.AddPostHandler("/api/recipe/upload_picture", recipePictureUploadHandler.Handle)
}
