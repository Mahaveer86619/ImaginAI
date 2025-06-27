package main

import (
	"fmt"
	"net/http"
	"os"

	postgres "github.com/Mahaveer86619/ImaginAI/src/database"
	handlers "github.com/Mahaveer86619/ImaginAI/src/handlers"
	middleware "github.com/Mahaveer86619/ImaginAI/src/middleware"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const defaultPort = "5050"

func main() {
	// Initialize logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	mux := http.NewServeMux()

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load environment variables")
	} else {
		logrus.Info("Environment variables loaded")
	}

	// Connect to database
	db, err := postgres.ConnectDB()
	if err != nil {
		logrus.WithError(err).Fatal("Error connecting to database")
	} else {
		logrus.Info("Database connected successfully")
	}
	defer postgres.CloseDBConnection(db)

	// Create tables if they don't exist
	err = postgres.CreateTables(db)
	if err != nil {
		logrus.WithError(err).Fatal("Error creating/verifying tables")
	} else {
		logrus.Info("Tables created/verified")
	}

	postgres.SetDBConnection(db)

	handleFunctions(mux)

	// Wrap all routes with the CORS and logging middleware
	handler := middleware.CORSMiddleware(middleware.LoggingMiddleware(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	addr := "0.0.0.0:" + port

	// Start server
	logrus.Infof("Starting server on port %s...", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		logrus.WithError(err).Fatal("Error running server")
	}
}

func handleFunctions(mux *http.ServeMux) {
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/favicon.ico" {
			return
		}
		fmt.Fprint(w, "ImaginAi API is running!")
	})

	//* Auth routes - POST methods
	mux.HandleFunc("/api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.RegisterUserController(w, r)
	})

	mux.HandleFunc("/api/v1/auth/authenticate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.AuthenticateUserController(w, r)
	})

	mux.HandleFunc("/api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.RefreshTokenController(w, r)
	})

	//* User routes - different methods for same path
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserByIDController(w, r)
		case http.MethodPut:
			handlers.UpdateUserController(w, r)
		case http.MethodDelete:
			handlers.DeleteUserController(w, r)
		case http.MethodOptions:
			return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	//* Admin routes - GET only
	mux.HandleFunc("/api/v1/users/all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.GetAllUsersController(w, r)
	})
}
