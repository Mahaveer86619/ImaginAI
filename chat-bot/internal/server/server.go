package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/cors"
	"google.golang.org/genai"
)

type GenAIServer struct {
	Ctx    context.Context
	client *genai.Client
	mu     sync.Mutex
}

func New(ctx context.Context) *GenAIServer {
	return &GenAIServer{
		Ctx: ctx,
		mu:  sync.Mutex{},
	}
}

func (s *GenAIServer) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/favicon.ico" {
			return
		}
		fmt.Fprint(w, "ImaginAi chat bot is running!")
	})
	mux.HandleFunc("/chat", s.ChatHandler)
	mux.HandleFunc("/stream", s.StreamChatHandler)
	mux.HandleFunc("/setup", s.SetupHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "Content-Type"},
	})

	return c.Handler(mux)
}

// // Wrapper for chatHandler to check Gemini initialization
// func (s *GenAIServer) chatHandlerWithCheck(w http.ResponseWriter, r *http.Request) {
// 	s.mu.Lock()
// 	initialized := *s.GeminiClient != nil && s.ModelName != ""
// 	s.mu.Unlock()
// 	if !initialized {
// 		http.Error(w, "Gemini API key not set. Please POST to /setup first.", http.StatusServiceUnavailable)
// 		return
// 	}
// 	s.chatHandler(w, r)
// }

// // Wrapper for streamingChatHandler to check Gemini initialization
// func (s *GenAIServer) streamingChatHandlerWithCheck(w http.ResponseWriter, r *http.Request) {
// 	s.mu.Lock()
// 	initialized := *s.GeminiClient != nil && s.ModelName != ""
// 	s.mu.Unlock()
// 	if !initialized {
// 		http.Error(w, "Gemini API key not set. Please POST to /setup first.", http.StatusServiceUnavailable)
// 		return
// 	}
// 	s.streamingChatHandler(w, r)
// }
