package services

import (
	"os"
	"sort"
	"strings"
)

type RAGChunk struct {
	ID     int
	Text   string
	TF     map[string]int
	Length int
}

type RAGService struct {
	chunks []RAGChunk
	df     map[string]int
	total  int
}

func NewRAGService(docPath string) (*RAGService, error) {
	content, err := os.ReadFile(docPath)
	if err != nil {
		return nil, err
	}

	text := normalizeText(string(content))
	chunks := splitIntoChunks(text, 900)

	service := &RAGService{
		chunks: make([]RAGChunk, 0, len(chunks)),
		df:     map[string]int{},
		total:  len(chunks),
	}

	for i, chunk := range chunks {
		tf := termFrequency(chunk)
		service.chunks = append(service.chunks, RAGChunk{
			ID:     i + 1,
			Text:   chunk,
			TF:     tf,
			Length: len(chunk),
		})
		for token := range tf {
			service.df[token]++
		}
	}

	return service, nil
}

func (r *RAGService) Retrieve(query string, k int) []RAGChunk {
	if k <= 0 || r.total == 0 {
		return nil
	}

	queryTF := termFrequency(query)
	type scored struct {
		chunk RAGChunk
		score float64
	}
	scoredChunks := make([]scored, 0, len(r.chunks))

	for _, chunk := range r.chunks {
		score := 0.0
		for token, qCount := range queryTF {
			df := r.df[token]
			if df == 0 {
				continue
			}
			idf := 1.0 + (float64(r.total) / float64(1+df))
			tf := chunk.TF[token]
			if tf == 0 {
				continue
			}
			score += float64(tf*qCount) * idf
		}
		if score > 0 {
			scoredChunks = append(scoredChunks, scored{chunk: chunk, score: score})
		}
	}

	sort.Slice(scoredChunks, func(i, j int) bool {
		return scoredChunks[i].score > scoredChunks[j].score
	})

	if len(scoredChunks) > k {
		scoredChunks = scoredChunks[:k]
	}

	results := make([]RAGChunk, 0, len(scoredChunks))
	for _, item := range scoredChunks {
		results = append(results, item.chunk)
	}
	return results
}

func normalizeText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.TrimSpace(text)
}

func splitIntoChunks(text string, maxLen int) []string {
	paragraphs := strings.Split(text, "\n\n")
	chunks := make([]string, 0, len(paragraphs))
	var current strings.Builder

	flush := func() {
		if current.Len() == 0 {
			return
		}
		chunk := strings.TrimSpace(current.String())
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		current.Reset()
	}

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if current.Len()+len(p)+2 > maxLen {
			flush()
		}
		if current.Len() > 0 {
			current.WriteString("\n\n")
		}
		current.WriteString(p)
	}
	flush()
	return chunks
}

func termFrequency(text string) map[string]int {
	tokens := tokenize(text)
	tf := make(map[string]int, len(tokens))
	for _, token := range tokens {
		if len(token) < 2 {
			continue
		}
		tf[token]++
	}
	return tf
}

func tokenize(text string) []string {
	text = strings.ToLower(text)
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		if r >= 'a' && r <= 'z' {
			return false
		}
		if r >= '0' && r <= '9' {
			return false
		}
		return true
	})
	return tokens
}
