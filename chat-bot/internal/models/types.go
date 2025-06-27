package models

import (
	"github.com/google/generative-ai-go/genai"
)

type Part struct {
	Text string
}

type Content struct {
	Role  string
	Parts []Part
}

type APIKeyRequest struct {
	APIKey string `json:"api_key"`
}

func (c *Content) Transform() *genai.Content {
	gc := &genai.Content{}
	gc.Role = c.Role
	ps := make([]genai.Part, len(c.Parts))
	for i, p := range c.Parts {
		ps[i] = genai.Text(p.Text)
	}
	gc.Parts = ps
	return gc
}

func Transform(cs []Content) []*genai.Content {
	gcs := make([]*genai.Content, len(cs))
	for i, c := range cs {
		gcs[i] = c.Transform()
	}
	return gcs
}
