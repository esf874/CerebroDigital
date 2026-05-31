package application

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// IDGenerator es un generador de IDs.
type IDGenerator func() string

// NoteService coordina la creación, eliminación, actualización y consulta de notas, gestionando también los enlaces asociados.
type NoteService struct {
	graph       *domain.Graph
	noteRepo    domain.InterfaceNoteRepository
	linkRepo    domain.InterfaceLinkRepository
	linkService *LinkService
	generateID  IDGenerator
}

func NewNoteService(
	graph *domain.Graph,
	noteRepo domain.InterfaceNoteRepository,
	linkRepo domain.InterfaceLinkRepository,
	linkService *LinkService,
	generateID IDGenerator,
) *NoteService {
	return &NoteService{
		graph:       graph,
		noteRepo:    noteRepo,
		linkRepo:    linkRepo,
		linkService: linkService,
		generateID:  generateID,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, title, content string) (*domain.Note, error) {

	// Generar objectID
	id := s.generateID()

	note, err := domain.NewNote(id, title, content)
	if err != nil {
		return nil, err
	}

	err = s.noteRepo.Save(ctx, note)
	if err != nil {
		return nil, err
	}

	// Actualizar grafo
	if err := s.graph.AddNote(note); err != nil {
		// Rollback, si no lo anhade en grafo, lo quita de base de datos
		_ = s.noteRepo.Delete(ctx, note.ID)
		return nil, err
	}

	// Procesar enlaces al crear
	titles := extractLinks(content)

	for _, t := range titles {
		var targetID string

		for _, n := range s.graph.Notes {
			if strings.EqualFold(n.Title, t) {
				targetID = n.ID
				break
			}
		}

		if targetID != "" {
			_, _ = s.linkService.CreateLink(ctx, note.ID, targetID, "")
		}
	}

	return note, nil
}

func (s *NoteService) GetNote(ctx context.Context, id string) (*domain.Note, error) {
	return s.graph.GetNote(id)
}

// Elimina nota y sus enlaces asociados
func (s *NoteService) DeleteNote(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	// Obtener y eliminar links asociados en la bd
	if err := s.linkRepo.DeleteAllByNoteID(ctx, id); err != nil {
		return fmt.Errorf("Failed to delete links: %w", err)
	}

	// Eliminar de la base de datos
	if err := s.noteRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("Failed to delete note: %w", err)
	}

	// Eliminar nota del grafo (con links asociados)
	return s.graph.RemoveNote(id)
}

func (s *NoteService) UpdateNote(ctx context.Context, id string, update domain.NoteUpdate) error {
	note, err := s.graph.GetNote(id)
	if err != nil {
		return err
	}

	// Copia de seguridad para rollback en caso de error
	tmp := *note

	// Aplicar cambios (tras validar)
	cambiado, err := note.ApplyUpdate(update)
	if err != nil {
		return err
	}

	if !cambiado {
		return nil
	}

	// Gestionar enlaces
	if update.Content != nil {
		// Extrae título
		titles := extractLinks(*update.Content)
		_ = s.linkRepo.DeleteAllByNoteID(ctx, id)
		s.graph.RemoveLinksByOrigin(id)

		// Convertir de títulos a ids
		for _, t := range titles {
			// Comprobar que existe nota con título
			var targetID string
			for _, n := range s.graph.Notes {
				if strings.EqualFold(n.Title, t) {
					targetID = n.ID
					break
				}
			}

			// Crea el enlace
			if targetID != "" {
				_, _ = s.linkService.CreateLink(ctx, id, targetID, "")
			}
		}
	}

	if err := s.noteRepo.Save(ctx, note); err != nil {
		*note = tmp // Rollback, restaurar como memoria
		return err
	}

	return nil
}

// Método auxiliar para extraer [[links]] del contenido
func extractLinks(content string) []string {
	re := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	matches := re.FindAllStringSubmatch(content, -1)

	ids := []string{}
	for _, match := range matches {
		if len(match) > 1 {
			ids = append(ids, match[1])
		}
	}
	return ids
}

func (s *NoteService) AddTagToNote(ctx context.Context, noteID, tag string) error {
	if tag == "" {
		return domain.ErrEmptyTag
	}

	note, err := s.graph.GetNote(noteID)
	if err != nil {
		return err
	}

	if err := note.AddTag(tag); err != nil {
		return err
	}
	if err := s.noteRepo.Save(ctx, note); err != nil {
		return err
	}

	return nil
}

func (s *NoteService) RemoveTagFromNote(ctx context.Context, noteID, tag string) error {
	if tag == "" {
		return domain.ErrEmptyTag
	}

	note, err := s.graph.GetNote(noteID)
	if err != nil {
		return err
	}

	if err := note.RemoveTag(tag); err != nil {
		return err
	}
	if err := s.noteRepo.Save(ctx, note); err != nil {
		return err
	}

	return nil
}

func (s *NoteService) Search(ctx context.Context, query string) ([]*domain.Note, error) {
	return s.noteRepo.SearchByContent(ctx, query)
}

func (s *NoteService) SearchByTag(ctx context.Context, tag string) ([]*domain.Note, error) {
	return s.noteRepo.FindByTag(ctx, tag)
}

func (s *NoteService) GetAllNotes(ctx context.Context) ([]*domain.Note, error) {
	notes := s.graph.GetAllNotes()

	// Ordenar segun fecha de actualización
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].UpdatedAt.After(notes[j].UpdatedAt)
	})
	return notes, nil
}
