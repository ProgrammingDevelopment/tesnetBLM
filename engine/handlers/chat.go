package handlers

import (
	"net/http"
	"strings"
	"war-ticket-engine/services"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func ChatHandler(rag *services.RAGService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		req.Message = strings.TrimSpace(req.Message)
		if req.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
			return
		}

		contexts := rag.Retrieve(req.Message, 4)
		answer, err := services.GenerateWithOllama(req.Message, contexts)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"answer":  "Model unavailable. Check OLLAMA_BASE_URL and OLLAMA_MODEL.",
				"sources": trimSources(contexts, 600),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"answer":  answer,
			"sources": trimSources(contexts, 600),
		})
	}
}

func trimSources(chunks []services.RAGChunk, maxLen int) []string {
	if len(chunks) == 0 {
		return nil
	}
	out := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		text := chunk.Text
		if len(text) > maxLen {
			text = text[:maxLen] + "..."
		}
		out = append(out, text)
	}
	return out
}
