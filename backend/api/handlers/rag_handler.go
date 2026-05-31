package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/application"
)

// AskRequest es el modelo para la petición.
type AskRequest struct {
	Question      string `json:"question"`
	CurrentNoteID string `json:"currentNoteId,omitempty"` // opcional
}

// AskResponse es el modelo de respuesta.
type AskResponse struct {
	Answer string `json:"answer"`
}

type RAGHandler struct {
	ragService *application.RAGService
}

func NewRAGHandler(ragService *application.RAGService) *RAGHandler {
	return &RAGHandler{ragService: ragService}
}

func (h *RAGHandler) Ask(w http.ResponseWriter, r *http.Request) {
	// Decodificar peticion JSON a struct go
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Question == "" {
		respondWithError(w, http.StatusBadRequest, "Question is required")
		return
	}

	answer, err := h.ragService.Ask(r.Context(), req.CurrentNoteID, req.Question)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AskResponse{Answer: answer})
}
