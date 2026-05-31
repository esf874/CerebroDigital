// backend/internal/api/handlers/graph_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/application"
)

type GraphHandler struct {
	graphService *application.GraphService
}

func NewGraphHandler(graphService *application.GraphService) *GraphHandler {
	return &GraphHandler{graphService: graphService}
}

// GET /api/notes/{id}/graph
func (h *GraphHandler) GetGraph(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	id := strings.TrimSuffix(path, "/graph")

	if id == "" {
		http.Error(w, `{"error": "ID requerido"}`, http.StatusBadRequest)
		return
	}

	depth := 2
	if d := r.URL.Query().Get("depth"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			depth = parsed
		}
	}

	limit := 30
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	graph, err := h.graphService.GetSubgraph(id, depth, limit)
	if err != nil {
		http.Error(w, `{"error": "Error obteniendo grafo"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(graph)
}
