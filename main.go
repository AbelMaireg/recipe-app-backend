package main

import (
	"log"
	"net/http"

	"app/config"
	"app/framework"
	"app/handlers"
	"app/repositories"
	"app/services"
)

func main() {
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

	router := framework.GetRouter()
	handlers.SetupRoutes(router, userService)

	addr := ":8080"
	log.Printf("Server running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, router.Instance))
}
