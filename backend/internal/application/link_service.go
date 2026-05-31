package application

import (
	"context"
	"fmt"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// LinkService coordina la creación, eliminación y consulta de enlaces entre notas.
type LinkService struct {
	graph      *domain.Graph
	linkRepo   domain.InterfaceLinkRepository
	noteRepo   domain.InterfaceNoteRepository
	generateID IDGenerator
}

func NewLinkService(
	graph *domain.Graph,
	linkRepo domain.InterfaceLinkRepository,
	noteRepo domain.InterfaceNoteRepository,
	generateID IDGenerator,
) *LinkService {
	return &LinkService{
		graph:      graph,
		linkRepo:   linkRepo,
		noteRepo:   noteRepo,
		generateID: generateID,
	}
}

func (s *LinkService) CreateLink(ctx context.Context, originID, destID, alias string) (*domain.Link, error) {

	if _, err := s.graph.GetNote(originID); err != nil {
		return nil, fmt.Errorf("Origin note validation failed: %w", err)
	}
	destNote, err := s.graph.GetNote(destID)
	if err != nil {
		return nil, fmt.Errorf("Destination note validation failed: %w", err)
	}

	linkID := s.generateID()

	link, err := domain.NewLink(linkID, originID, destID, alias, destNote.Title)
	if err != nil {
		return nil, err
	}

	err = s.linkRepo.Save(ctx, link)
	if err != nil {
		return nil, err
	}

	if err := s.graph.AddLink(link); err != nil {
		return nil, fmt.Errorf("failed to add link to graph: %w", err)
	}
	return link, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, id string) error {
	// Eliminar de la bd
	if err := s.linkRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("Delete link: %w", err)
	}

	// Quitar del grafo
	if err := s.graph.RemoveLink(id); err != nil {
		if err == domain.ErrLinkNotFound {
			return nil
		}
		return err
	}
	return nil
}

// Consultas por links
func (s *LinkService) GetByOrigin(ctx context.Context, originID string) ([]*domain.Link, error) {
	return s.graph.OutgoingLinks(originID)
}

func (s *LinkService) GetByDest(ctx context.Context, destID string) ([]*domain.Link, error) {
	return s.graph.IncomingLinks(destID)
}

func (s *LinkService) GetAllForNote(ctx context.Context, noteID string) ([]*domain.Link, error) {
	return s.graph.AllLinksForNote(noteID)
}
