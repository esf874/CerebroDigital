package domain

import (
	"context"
)

type InterfaceNoteRepository interface {
	Save(ctx context.Context, note *Note) error // upsert
	FindById(ctx context.Context, id string) (*Note, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*Note, error)
	FindByTag(ctx context.Context, tag string) ([]*Note, error)
	SearchByContent(ctx context.Context, query string) ([]*Note, error)
}

type InterfaceLinkRepository interface {
	Save(ctx context.Context, link *Link) error
	Delete(ctx context.Context, id string) error
	DeleteAllByNoteID(ctx context.Context, noteID string) error
	FindByOrigin(ctx context.Context, noteId string) ([]*Link, error)
	FindByDest(ctx context.Context, noteId string) ([]*Link, error)
	FindByNoteID(ctx context.Context, noteId string) ([]*Link, error)
	FindAll(ctx context.Context) ([]*Link, error)
}
