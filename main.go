package main

import (
	"log"
	"net/http"

	"app/framework"
	"app/handlers"
)

func main() {
	handlers.SetupRoutes(framework.GetRouter())

	addr := ":8080"
	log.Printf("Server running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, framework.GetRouter().Instance))
}
