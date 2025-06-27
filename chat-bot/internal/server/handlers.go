package server

import (
	"encoding/json"
	"net/http"

	"github.com/Mahaveer86619/ImaginAI/internal/models"
	"google.golang.org/genai"
)

const (
	RoleUser  string = genai.RoleUser
	RoleModel string = genai.RoleModel
)

func (gs *GenAIServer) ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gs.mu.Lock()
	client := gs.client
	gs.mu.Unlock()

	if client == nil {
		http.Error(w, "Gemini client not initialized. Please POST /setup first.", http.StatusServiceUnavailable)
		return
	}

	temp := float32(0.9)
	topP := float32(0.5)
	topK := float32(20.0)

	config := &genai.GenerateContentConfig{
		Temperature: &temp,
		TopP:        &topP,
		TopK:        &topK,
	}

	history := req.History
	history = append(history, models.ChatMessage{Role: RoleUser, Message: req.Message})

	chat, err := client.Chats.Create(gs.Ctx, "gemini-2.5-flash", config, ToGenaiContent(history))
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	res, err := chat.SendMessage(gs.Ctx, genai.Part{Text: req.Message})
	if err != nil || res == nil {
		http.Error(w, "Failed to send message to Gemini. err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	history = append(history, models.ChatMessage{Role: RoleModel, Message: res.Text()})

	resp := models.ChatResponse{
		Response: res.Text(),
		History:  history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (gs *GenAIServer) StreamChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gs.mu.Lock()
	client := gs.client
	gs.mu.Unlock()

	if client == nil {
		http.Error(w, "Gemini client not initialized. Please POST /setup first.", http.StatusServiceUnavailable)
		return
	}

	temp := float32(0.9)
	topP := float32(0.5)
	topK := float32(20.0)

	config := &genai.GenerateContentConfig{
		Temperature: &temp,
		TopP:        &topP,
		TopK:        &topK,
	}

	history := req.History
	history = append(history, models.ChatMessage{Role: RoleUser, Message: req.Message})

	chat, err := client.Chats.Create(gs.Ctx, "gemini-2.5-flash", config, ToGenaiContent(history))
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	iter := chat.SendMessageStream(gs.Ctx, genai.Part{Text: req.Message})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	var responseText string

	for chunk, _ := range iter {
		if chunk == nil {
			break
		}
		part := chunk.Candidates[0].Content.Parts[0]
		responseText += part.Text
		// Send each part as an SSE event
		_, _ = w.Write([]byte("data: " + responseText + "\n\n"))
		flusher.Flush()
	}

	// Save the model's response to history
	history = append(history, models.ChatMessage{Role: RoleModel, Message: responseText})

	// Optionally, send the final history as a JSON event at the end
	finalResp := models.ChatResponse{
		Response: responseText,
		History:  history,
	}
	finalJSON, _ := json.Marshal(finalResp)
	_, _ = w.Write([]byte("event: history\ndata: " + string(finalJSON) + "\n\n"))
	flusher.Flush()
}

func (s *GenAIServer) SetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.APIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.APIKey == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	client, err := genai.NewClient(s.Ctx, &genai.ClientConfig{
		APIKey: req.APIKey,
	})
	if err != nil {
		http.Error(w, "Failed to initialize Gemini client", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	s.client = client
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ready"}`))
}

func ToGenaiContent(history []models.ChatMessage) []*genai.Content {
	history_content := []*genai.Content{}
	for _, msg := range history {
		var role genai.Role
		if msg.Role == RoleUser {
			role = genai.RoleUser
		} else {
			role = genai.RoleModel
		}

		content := genai.NewContentFromText(msg.Message, role)
		history_content = append(history_content, content)
	}
	return history_content
}
