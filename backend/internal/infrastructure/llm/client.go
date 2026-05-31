package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message representa un mensaje con formato estandar de la API de OpenAI.
type Message struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

// ChatRequest representa la estructura de la petición a llama-server.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

// ChatResponse representa la respuesta de llama-server.
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

// Adaptador para comunicacion con llama-server
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Chat envía un prompt al LLM y devuelve la respuesta procesada, manejando limpieza de texto y casos sin información.
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {

	messages := []Message{
		{
			Role: "system",
			Content: `Eres un asistente de gestión de conocimiento personal.
			REGLAS OBLIGATORIAS:
			1. Responde SIEMPRE en español.
			2. Responde de forma DIRECTA y CONCISA.
			3. Si la respuesta está en las notas, dala directamente.
			4. Si no hay información en las notas, responde: "No tengo información sobre eso en mis notas."
			5. Para listas de recetas o elementos, usa formato de lista simple con guiones (-).`,
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	reqBody := ChatRequest{
		Model:       "qwen",
		Messages:    messages,
		Stream:      false,
		Temperature: 0.5,
		MaxTokens:   1024,
	}

	// Convertir peticion a JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("Error marshaling request: %w", err)
	}

	// Crear petición HTTP
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+"/v1/chat/completions", // endpoint estandar de OpenAI
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ejecutar la petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error calling llama-server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Llama-server returned status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return "", fmt.Errorf("Error decoding response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("No choices in response")
	}

	// Extraer y limpiar respuesta del modelo, si es muy larga o contiene razonamiento, eliminar lo innecesario
	msg := chatResp.Choices[0].Message
	finalAnswer := msg.Content

	if finalAnswer == "" {
		finalAnswer = msg.ReasoningContent
	}

	// Limpieza de etiquetas de razonamiento
	removePhrases := []string{
		"Okay, let me",
		"Let me first",
		"I need to",
		"First, I'll",
		"Wait,",
		"Let me check",
		"Based on the",
		"The user is asking",
		"I'll structure",
		"Let me try to",
		"RESPUESTA:",
		"Answer:",
	}

	for _, phrase := range removePhrases {
		if idx := strings.Index(finalAnswer, phrase); idx != -1 {
			if idx+len(phrase) < len(finalAnswer) {
				finalAnswer = finalAnswer[idx+len(phrase):]
			}
		}
	}

	lines := strings.Split(finalAnswer, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "===") {
			continue
		}
		if strings.Contains(line, "thought") || strings.Contains(line, "reasoning") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	finalAnswer = strings.Join(cleanLines, "\n")
	finalAnswer = strings.TrimSpace(finalAnswer)

	return finalAnswer, nil
}
