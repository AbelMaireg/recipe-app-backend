package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-graphql-app/models"
	"go-graphql-app/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type SignInResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
}

type HasuraAction struct {
	Action struct {
		Name string `json:"name"`
	} `json:"action"`
	Input json.RawMessage `json:"input"`
}

type UserEventPayload struct {
	Event struct {
		Op   string `json:"op"`
		Data struct {
			New models.User `json:"new"`
		} `json:"data"`
	} `json:"event"`
}

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=secret dbname=userapp port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	// Set up HTTP router
	r := chi.NewRouter()

	// Hasura action handler
	r.Post("/actions", func(w http.ResponseWriter, r *http.Request) {
		var action HasuraAction
		if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		switch action.Action.Name {
		case "signUp":
			var input SignUpInput
			if err := json.Unmarshal(action.Input, &input); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			hashedPassword, err := utils.HashPassword(input.Password)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}

			dbUser := models.User{
				Username: input.Username,
				Password: hashedPassword,
			}
			if err := db.Create(&dbUser).Error; err != nil {
				http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusBadRequest)
				return
			}

			response := SignUpResponse{
				ID:       fmt.Sprintf("%d", dbUser.ID),
				Username: dbUser.Username,
			}
			json.NewEncoder(w).Encode(response)

		case "signIn":
			var input SignInInput
			if err := json.Unmarshal(action.Input, &input); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			var dbUser models.User
			if err := db.Where("username = ?", input.Username).First(&dbUser).Error; err != nil {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			if err := utils.VerifyPassword(dbUser.Password, input.Password); err != nil {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			token, err := utils.GenerateJWT(dbUser.ID)
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}

			response := SignInResponse{
				Token: token,
			}
			response.User.ID = fmt.Sprintf("%d", dbUser.ID)
			response.User.Username = dbUser.Username
			json.NewEncoder(w).Encode(response)

		default:
			http.Error(w, "Unknown action", http.StatusBadRequest)
		}
	})

	// Hasura event handler
	r.Post("/events", func(w http.ResponseWriter, r *http.Request) {
		var payload UserEventPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid event payload", http.StatusBadRequest)
			return
		}

		if payload.Event.Op == "INSERT" {
			log.Printf("User created: ID=%d, Username=%s", payload.Event.Data.New.ID, payload.Event.Data.New.Username)
			// In a real app, insert into audit_log table or send notification
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Println("Action server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
