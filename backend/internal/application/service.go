package application

import (
	"context"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// NoteServiceInterface define la interfaz pública de NoteService, facilitando el desacoplamiento y testing.
type NoteServiceInterface interface {
	CreateNote(ctx context.Context, title, content string) (*domain.Note, error)
	GetNote(ctx context.Context, id string) (*domain.Note, error)
	UpdateNote(ctx context.Context, id string, update domain.NoteUpdate) error
	DeleteNote(ctx context.Context, id string) error
	Search(ctx context.Context, query string) ([]*domain.Note, error)
	SearchByTag(ctx context.Context, tag string) ([]*domain.Note, error)
	GetAllNotes(ctx context.Context) ([]*domain.Note, error)
	RemoveTagFromNote(ctx context.Context, noteID, tag string) error
	AddTagToNote(ctx context.Context, noteID, tag string) error
}

// LinkServiceInterface define la interfaz pública de LinkService, facilitando el desacoplamiento y testing.
type LinkServiceInterface interface {
	CreateLink(ctx context.Context, originID, destID, alias string) (*domain.Link, error)
	DeleteLink(ctx context.Context, id string) error
	GetByOrigin(ctx context.Context, originID string) ([]*domain.Link, error)
	GetByDest(ctx context.Context, destID string) ([]*domain.Link, error)
	GetAllForNote(ctx context.Context, noteID string) ([]*domain.Link, error)
}

// GraphServiceInterface define la interfaz pública de GraphService, facilitando el desacoplamiento y testing.
type GraphServiceInterface interface {
	LoadGraph(ctx context.Context) error
	GetConnectedNotes(noteID string) ([]*domain.Note, error)
	OutgoingLinks(noteID string) ([]*domain.Link, error)
	IncomingLinks(noteID string) ([]*domain.Link, error)
	GetAllLinks(noteID string) ([]*domain.Link, error)
	GetSubgraph(noteID string, depth int, limit int) (*GraphResponse, error)
}

// RAGServiceInterface define la interfaz pública de RAGService, facilitando el desacoplamiento y testing.
type RAGServiceInterface interface {
	Ask(ctx context.Context, currentNoteID string, userQuery string) (string, error)
	AskWithContext(ctx context.Context, currentNoteID string, userQuery string) (string, error)
	exploreGraph(startNoteID string, maxDepth int) ([]*domain.Note, error)
	expandByTags(tags []string, queue *[]string, visited map[string]bool)
	buildPromptWithCurrentNote(query string, notes []*domain.Note, currentNote *domain.Note) string
	formatNote(note *domain.Note) string
	AskGlobal(ctx context.Context, userQuery string) (string, error)
	findRelevantNotes(ctx context.Context, query string) ([]*domain.Note, error)
	buildPromptWithNotes(query string, notes []*domain.Note) string
	askWithoutContext(ctx context.Context, query string) (string, error)
}
