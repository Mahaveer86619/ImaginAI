package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Mahaveer86619/ImaginAI/internal/server"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort = "5000"
)

func main() {
	// Initialize logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	ctx := context.Background()
	srv := server.New(ctx) // This instance must be reused for all requests

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	addr := "0.0.0.0:" + port

	// Start server
	logrus.Infof("Starting chat bot on port %s...", port)
	if err := http.ListenAndServe(addr, srv.SetupRoutes()); err != nil {
		logrus.WithError(err).Fatal("Error running chat bot")
	}
}

// func initGeminiClient(apiKey string) {
// 	ctx := context.Background()
// 	fmt.Print("Creating Gemini client...", apiKey)
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
// 	if err != nil {
// 		log.Fatalf("could not create Gemini client. err: %v", err)
// 	}
// 	GeminiClient = client
// 	GeminiModel = client.GenerativeModel(modelName)
// }
