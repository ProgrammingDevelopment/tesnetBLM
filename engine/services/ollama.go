package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string       `json:"model"`
	Messages []LLMMessage `json:"messages"`
	Stream   bool         `json:"stream"`
}

type ollamaChatResponse struct {
	Message LLMMessage `json:"message"`
}

func GenerateWithOllama(question string, contexts []RAGChunk) (string, error) {
	baseURL := strings.TrimSpace(os.Getenv("OLLAMA_BASE_URL"))
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	model := strings.TrimSpace(os.Getenv("OLLAMA_MODEL"))
	if model == "" {
		model = "llama3"
	}

	contextText := buildContextText(contexts)
	systemPrompt := "You are a helpful assistant for War Tiket Engine. " +
		"Answer using the provided context. If the answer is not in the context, say you do not have that information."
	userPrompt := "Context:\n" + contextText + "\n\nQuestion:\n" + question

	reqBody := ollamaChatRequest{
		Model: model,
		Messages: []LLMMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseURL, "/")+"/api/chat", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("ollama request failed")
	}

	var out ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	answer := strings.TrimSpace(out.Message.Content)
	if answer == "" {
		return "", errors.New("empty model response")
	}
	return answer, nil
}

func buildContextText(contexts []RAGChunk) string {
	if len(contexts) == 0 {
		return "No relevant context found."
	}
	var b strings.Builder
	for i, ctx := range contexts {
		b.WriteString("Source ")
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(":\n")
		b.WriteString(ctx.Text)
		b.WriteString("\n\n")
	}
	return strings.TrimSpace(b.String())
}
