package models

type ChatMessage struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type ChatRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history"`
}

type ChatResponse struct {
	Response string        `json:"response"`
	History  []ChatMessage `json:"history"`
}


