package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/application"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// CreateNoteRequest es el formato para crear una nota.
type CreateNoteRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
	Priority string   `json:"priority"`
}

// UpdateNoteRequest es el formato para actualizar una nota.
type UpdateNoteRequest struct {
	Title    *string   `json:"title,omitempty"`
	Content  *string   `json:"content,omitempty"`
	Theme    *string   `json:"theme,omitempty"`
	Tags     *[]string `json:"tags,omitempty"`
	Status   *string   `json:"status,omitempty"`
	Priority *string   `json:"priority,omitempty"`
}

// NoteResponse es el formato de respuesta para una nota.
type NoteResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Theme     string   `json:"theme,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Status    string   `json:"status,omitempty"`
	Priority  string   `json:"priority,omitempty"`
	UpdatedAt string   `json:"updatedAt"`
}

// NoteHandler se encarga de manejar las peticiones relacionadas con notas.
type NoteHandler struct {
	noteService *application.NoteService
}

func NewNoteHandler(noteService *application.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

// GET /api/notes
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.GetAllNotes(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error interno")
		return
	}

	response := make([]NoteResponse, len(notes))
	for i, note := range notes {
		response[i] = toNoteResponse(note)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/notes/{id}
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	// extraer id de la url para obtener la nota
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID requerido")
		return
	}

	note, err := h.noteService.GetNote(r.Context(), id)
	if err != nil {
		if err == domain.ErrNoteNotFound {
			respondWithError(w, http.StatusNotFound, "Nota no encontrada")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error interno")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toNoteResponse(note))
}

func (h *NoteHandler) GetNoteByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		respondWithError(w, http.StatusBadRequest, "Título necesario")
		return
	}
	notes, _ := h.noteService.GetAllNotes(r.Context())
	for _, n := range notes {
		if strings.EqualFold(n.Title, title) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(toNoteResponse(n))
			return
		}
	}
	respondWithError(w, http.StatusNotFound, "Nota no encontrada")
}

// POST /api/notes
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.Title == "" {
		respondWithError(w, http.StatusBadRequest, "El título es requerido")
		return
	}

	note, err := h.noteService.CreateNote(r.Context(), req.Title, req.Content)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al crear nota")
		return
	}

	update := domain.NoteUpdate{}

	if req.Status != "" {
		status := domain.NoteStatus(req.Status)
		update.Status = &status
	}
	if req.Priority != "" {
		priority := domain.NotePriority(req.Priority)
		update.Priority = &priority
	}

	if len(req.Tags) > 0 {
		update.Tags = &req.Tags
	}

	if _, err := note.ApplyUpdate(update); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al aplicar cambios")
		return
	}

	if err := h.noteService.UpdateNote(r.Context(), note.ID, update); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al guardar cambios")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toNoteResponse(note))
}

// PUT /api/notes/{id}
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID requerido")
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Construir nota actualizada a partir de la petición.
	update := domain.NoteUpdate{
		Title:   req.Title,
		Content: req.Content,
		Theme:   req.Theme,
		Tags:    req.Tags,
	}

	if req.Status != nil {
		status := domain.NoteStatus(*req.Status)
		update.Status = &status
	}
	if req.Priority != nil {
		priority := domain.NotePriority(*req.Priority)
		update.Priority = &priority
	}

	if err := h.noteService.UpdateNote(r.Context(), id, update); err != nil {
		if err == domain.ErrNoteNotFound {
			respondWithError(w, http.StatusNotFound, "Nota no encontrada")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar")
		return
	}

	updatedNote, _ := h.noteService.GetNote(r.Context(), id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toNoteResponse(updatedNote))
}

// DELETE /api/notes/{id}
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID requerido")
		return
	}

	if err := h.noteService.DeleteNote(r.Context(), id); err != nil {
		if err == domain.ErrNoteNotFound {
			respondWithError(w, http.StatusNotFound, "Nota no encontrada")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error al eliminar")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toNoteResponse(note *domain.Note) NoteResponse {
	return NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		Theme:     note.Theme,
		Tags:      note.Tags,
		Status:    string(note.Status),
		Priority:  string(note.Priority),
		UpdatedAt: note.UpdatedAt.Format(time.RFC3339), // Formato ISO
	}
}

// POST /api/notes/{id}/tags
func (h *NoteHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	parts := strings.Split(path, "/tags/")
	if len(parts) != 2 {
		respondWithError(w, http.StatusBadRequest, "URL inválida")
		return
	}

	vars := mux.Vars(r)
	noteID := vars["id"]
	tag := vars["tag"]

	if err := h.noteService.AddTagToNote(r.Context(), noteID, tag); err != nil {
		if err == domain.ErrTagAlreadyExists {
			respondWithError(w, http.StatusConflict, "La etiqueta ya existe")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error al añadir etiqueta")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/notes/{id}/tags/{tag}
func (h *NoteHandler) RemoveTag(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	parts := strings.Split(path, "/tags/")
	if len(parts) != 2 {
		respondWithError(w, http.StatusBadRequest, "URL inválida")
		return
	}

	noteID := parts[0]
	tag := parts[1]

	if err := h.noteService.RemoveTagFromNote(r.Context(), noteID, tag); err != nil {
		if err == domain.ErrNoteNotFound {
			respondWithError(w, http.StatusNotFound, "Nota no encontrada")
			return
		}
		if err == domain.ErrTagNotFound {
			respondWithError(w, http.StatusNotFound, "Etiqueta no encontrada")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error al eliminar etiqueta")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Error en formato JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
