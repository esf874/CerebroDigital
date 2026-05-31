package application

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// LLMClient es una interfaz, permitiendo la integrar cualquier modelo de lenguaje que implemente este método.
type LLMClient interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// RAGService coordina la lógica de recuperación de información relevante del grafo 
// y la interacción con el LLM para responder a las preguntas del usuario.
type RAGService struct {
	graph        *domain.Graph
	graphService *GraphService
	noteRepo     domain.InterfaceNoteRepository
	llmClient    LLMClient
	maxDepth     int
	maxNotes     int
}

func NewRAGService(graph *domain.Graph, graphService *GraphService, noteRepo domain.InterfaceNoteRepository, llmClient LLMClient) *RAGService {
	return &RAGService{
		graph:        graph,
		graphService: graphService,
		noteRepo:     noteRepo,
		llmClient:    llmClient,
		maxDepth:     2,
		maxNotes:     12,
	}
}

// Método para manejo de ambos casos: contextual y global.
func (s *RAGService) Ask(ctx context.Context, currentNoteID string, userQuery string) (string, error) {
	if currentNoteID != "" {
		return s.AskWithContext(ctx, currentNoteID, userQuery)
	}
	return s.AskGlobal(ctx, userQuery)
}

// Búsqueda contextual: explora el grafo desde la nota actual, recupera notas relacionadas, 
// construye un prompt con contexto y pregunta al LLM
func (s *RAGService) AskWithContext(ctx context.Context, currentNoteID string, userQuery string) (string, error) {

	currentNote, err := s.graph.GetNote(currentNoteID)
	if err != nil {
		return "", fmt.Errorf("Nota origen no encontrada: %w", err)
	}

	// Explorar el grafo desde la nota actual (BFS + tags)
	contextNotes, err := s.exploreGraph(currentNoteID, s.maxDepth)
	if err != nil {
		return "", fmt.Errorf("error explorando grafo: %w", err)
	}

	// La nota actual se coloca primero (la más relevante)
	allNotes := append([]*domain.Note{currentNote}, contextNotes...)

	// Priorizar la nota actual y las más cercanas, luego recortar si excede el límite
	if len(allNotes) > s.maxNotes {
		allNotes = allNotes[:s.maxNotes]
	}

	prompt := s.buildPromptWithCurrentNote(userQuery, allNotes, currentNote)

	response, err := s.llmClient.Chat(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("Error asking LLM: %w", err)
	}

	return response, nil
}

// BFS desde nota incial incluyendo expansion por tags
func (s *RAGService) exploreGraph(startNoteID string, maxDepth int) ([]*domain.Note, error) {
	
	result, err := s.graphService.bfs(startNoteID, maxDepth, 100)
	if err != nil {
		return nil, err
	}

	// Convertir IDs a notas
	notes := make([]*domain.Note, 0, len(result.Nodes))
	for _, n := range result.Nodes {
		if n.id == startNoteID {
			continue 
		}
		note, err := s.graph.GetNote(n.id)
		if err != nil {
			continue
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// buildPromptWithCurrentNote construye el prompt destacando la nota origen
func (s *RAGService) buildPromptWithCurrentNote(query string, notes []*domain.Note, currentNote *domain.Note) string {
	var builder strings.Builder

	// Instrucciones
	builder.WriteString(`Responde SIEMPRE en ESPAÑOL.
	Eres un asistente experto que responde exclusivamente basándose en las notas proporcionadas.
	REGLAS:
	1. Responde SOLO con la información de las notas.
	2. Si la información no está, di: "No dispongo de información suficiente en tus notas".
	3. Responde en español.

	NOTAS RELEVANTES:
	`)

	builder.WriteString(fmt.Sprintf("**ACTUAL:** %s\n", currentNote.Title))
	builder.WriteString(fmt.Sprintf("Contenido: %s\n\n", truncateString(currentNote.Content, 300)))

	// Si hay notas relacionadas, las listamos, si no, se responde sin contexto 
	if len(notes) > 1 {
		builder.WriteString("**RELACIONADAS:**\n")
		count := 0
		for _, note := range notes[1:] {
			if count >= 5 {
				break
			}
			builder.WriteString(fmt.Sprintf("- %s: %s\n", note.Title, truncateString(note.Content, 200)))
			count++
		}
		builder.WriteString("\n")
	}

	builder.WriteString(fmt.Sprintf("Pregunta: %s\n\nRespuesta:", query))
	return builder.String()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Convierte nota en texto formateado para el modelo, destacando título, tema, contenido y tags
func (s *RAGService) formatNote(note *domain.Note) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Título: %s\n", note.Title))
	if note.Theme != "" {
		b.WriteString(fmt.Sprintf("Tema: %s\n", note.Theme))
	}
	b.WriteString(fmt.Sprintf("Contenido: %s\n", note.Content))
	if len(note.Tags) > 0 {
		b.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(note.Tags, ", ")))
	}

	return b.String()
}

func (s *RAGService) AskGlobal(ctx context.Context, userQuery string) (string, error) {

	// Búsqueda notas relacionadas con la pregunta del usuario
	relevantNotes, err := s.findRelevantNotes(ctx, userQuery)
	if err != nil {
		return "", fmt.Errorf("Error finding relevant notes: %w", err)
	}

	// Si no hay notas, responder sin contexto
	if len(relevantNotes) == 0 {
		return s.askWithoutContext(ctx, userQuery)
	}

	prompt := s.buildPromptWithNotes(userQuery, relevantNotes)

	// Enviar al LLM con el prompt
	response, err := s.llmClient.Chat(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("Error asking LLM: %w", err)
	}

	return response, nil
}

// Búsqueda global por texto y tags optimizado
func (s *RAGService) findRelevantNotes(ctx context.Context, query string) ([]*domain.Note, error) {
	textMatches, err := s.noteRepo.SearchByContent(ctx, query)
	if err != nil {
		return nil, err
	}

	// Procesar keywords
	keywordsMap := make(map[string]bool)
	for _, k := range strings.Fields(strings.ToLower(query)) {
		keywordsMap[k] = true
	}

	// Ranking notas en base a posición en resultados de texto y coincidencia de tags
	type scoredNote struct {
		note  *domain.Note
		score float64
	}

	var scored []scoredNote
	for i, note := range textMatches {
		// Ranking base: posición + score semántico
		score := 100.0 / float64(i+1)

		for _, tag := range note.Tags {
			if keywordsMap[strings.ToLower(tag)] {
				score += 25.0 // Adicional por coincidencia de tag
			}
		}
		scored = append(scored, scoredNote{note, score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Selección mejores notas
	result := make([]*domain.Note, 0)
	totalChars := 0
	maxChars := 8000 // Límite para evitar overflow de tokens si notas muy extensas

	for _, sn := range scored {
		noteSize := len(sn.note.Content)
		if totalChars+noteSize > maxChars && len(result) > 0 {
			break // para si sig nota excede el contexto
		}
		result = append(result, sn.note)
		totalChars += noteSize
		if len(result) >= s.maxNotes {
			break
		}
	}
	return result, nil
}

// Construye el prompt enviado al LLM cuando no hay nota origen
func (s *RAGService) buildPromptWithNotes(query string, notes []*domain.Note) string {
	var builder strings.Builder

	builder.WriteString(`Responde SIEMPRE en ESPAÑOL.
	Eres un asistente experto que responde exclusivamente basándose en las notas proporcionadas.
	REGLAS:
	1. Responde SOLO con la información de las notas.
	2. Si la información no está, di: "No dispongo de información suficiente en tus notas".
	3. Responde en español.

	NOTAS RELEVANTES:
	`)

	for i, note := range notes {
		if i >= 5 {
			break
		}
		builder.WriteString(fmt.Sprintf("- %s: %s\n", note.Title, truncateString(note.Content, 200)))
	}

	builder.WriteString(fmt.Sprintf("\nPregunta: %s\n\nRespuesta:", query))
	return builder.String()
}

func (s *RAGService) askWithoutContext(ctx context.Context, query string) (string, error) {
	prompt := fmt.Sprintf(`Eres un asistente útil. El usuario pregunta: %s

	No se encontró información relevante en la base de conocimiento.
	Responde indicando que no tienes información específica sobre este tema.
	Respuesta:`, query)

	return s.llmClient.Chat(ctx, prompt)
}
