package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

func ParseRequestJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func RenderResponseJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func ResponseString(res *genai.GenerateContentResponse) (string, error) {
	if len(res.Candidates) > 0 {
		if cs := ContentString(res.Candidates[0].Content); cs != nil {
			return *cs, nil
		}
	}
	return "", fmt.Errorf("invalid response from Gemini model")
}

func ContentString(c *genai.Content) *string {
	if c == nil || c.Parts == nil {
		return nil
	}

	cStrs := make([]string, len(c.Parts))
	for i, part := range c.Parts {
		if pt, ok := part.(genai.Text); ok {
			cStrs[i] = string(pt)
		} else {
			return nil
		}
	}

	cStr := strings.Join(cStrs, "\n")
	return &cStr
}
