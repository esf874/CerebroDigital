package application

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// MockNoteRepository es un mock de InterfaceNoteRepository para tests, usando testify/mock.
type MockNoteRepository struct {
	mock.Mock 
}

func (m *MockNoteRepository) Save(ctx context.Context, note *domain.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockNoteRepository) FindById(ctx context.Context, id string) (*domain.Note, error) {
	args := m.Called(ctx, id)
	n, _ := args.Get(0).(*domain.Note)
	return n, args.Error(1)
}

func (m *MockNoteRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNoteRepository) FindAll(ctx context.Context) ([]*domain.Note, error) {
	args := m.Called(ctx)
	notes, _ := args.Get(0).([]*domain.Note)
	return notes, args.Error(1)
}

func (m *MockNoteRepository) FindByTag(ctx context.Context, tag string) ([]*domain.Note, error) {
	args := m.Called(ctx, tag)
	return args.Get(0).([]*domain.Note), args.Error(1)
}

func (m *MockNoteRepository) SearchByContent(ctx context.Context, query string) ([]*domain.Note, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*domain.Note), args.Error(1)
}

// MockLinkRepository es un mock de InterfaceLinkRepository para tests, usando testify/mock.
type MockLinkRepository struct {
	mock.Mock // anotacion @Mock
}

func (m *MockLinkRepository) Save(ctx context.Context, link *domain.Link) error {
	args := m.Called(ctx, link)
	return args.Error(0)
}

func (m *MockLinkRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLinkRepository) DeleteAllByNoteID(ctx context.Context, noteID string) error {
	args := m.Called(ctx, noteID)
	return args.Error(0)
}

func (m *MockLinkRepository) FindByOrigin(ctx context.Context, noteId string) ([]*domain.Link, error) {
	args := m.Called(ctx, noteId)
	links, _ := args.Get(0).([]*domain.Link)
	return links, args.Error(1)
}

func (m *MockLinkRepository) FindByDest(ctx context.Context, noteId string) ([]*domain.Link, error) {
	args := m.Called(ctx, noteId)
	links, _ := args.Get(0).([]*domain.Link)
	return links, args.Error(1)
}

func (m *MockLinkRepository) FindByNoteID(ctx context.Context, noteID string) ([]*domain.Link, error) {
	args := m.Called(ctx, noteID)
	return args.Get(0).([]*domain.Link), args.Error(1)
}

func (m *MockLinkRepository) FindAll(ctx context.Context) ([]*domain.Link, error) {
	args := m.Called(ctx)
	links, _ := args.Get(0).([]*domain.Link)
	return links, args.Error(1)
}
