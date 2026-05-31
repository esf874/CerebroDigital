package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// graphSvcMock es un helper para tests de GraphService, con mocks de repositorios y un grafo real.
type graphSvcMock struct {
	svc      *GraphService
	noteRepo *MockNoteRepository
	linkRepo *MockLinkRepository
	graph    *domain.Graph
}

func newGraphSvcMock(t *testing.T) *graphSvcMock {
	t.Helper()

	noteRepo := new(MockNoteRepository)
	linkRepo := new(MockLinkRepository)
	graph := domain.NewGraph()

	svc := NewGraphService(graph, noteRepo, linkRepo)

	fx := &graphSvcMock{
		svc:      svc,
		noteRepo: noteRepo,
		linkRepo: linkRepo,
		graph:    graph,
	}

	t.Cleanup(func() {
		noteRepo.AssertExpectations(t)
		linkRepo.AssertExpectations(t)
	})

	return fx
}

func TestGraphService_LoadGraph(t *testing.T) {

	t.Run("Éxito - carga notas y links para reconstruir grafo", func(t *testing.T) {
		fx := newGraphSvcMock(t)
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "N1", "C")
		n2, _ := domain.NewNote("n2", "N2", "C")
		notes := []*domain.Note{n1, n2}

		l1, _ := domain.NewLink("l1", "n1", "n2", "alias", "title")
		links := []*domain.Link{l1}

		// Se debe recuperar todas las notas y links
		fx.noteRepo.On("FindAll", mock.Anything).Return(notes, nil).Once()
		fx.linkRepo.On("FindAll", mock.Anything).Return(links, nil).Once()

		err := fx.svc.LoadGraph(ctx)

		require.NoError(t, err)
		_, err = fx.graph.GetNote("n1")
		require.NoError(t, err)

		_, err = fx.graph.GetNote("n2")
		require.NoError(t, err)

		gotlink, err := fx.graph.GetLink("l1")
		require.NoError(t, err)
		assert.Equal(t, "n1", gotlink.OriginNoteId)
		assert.Equal(t, "n2", gotlink.DestNoteId)
	})

	t.Run("Error - fallo del note repository", func(t *testing.T) {
		fx := newGraphSvcMock(t)
		ctx := context.Background()

		loadErr := errors.New("db notes down")
		fx.noteRepo.
			On("FindAll", mock.Anything).
			Return(([]*domain.Note)(nil), loadErr).Once()

		err := fx.svc.LoadGraph(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, loadErr)

		// No llama a findall links si fallan las notas
		fx.linkRepo.AssertNotCalled(t, "FindAll", mock.Anything)
	})

	t.Run("Error - fallo del link repositorio", func(t *testing.T) {
		fx := newGraphSvcMock(t)
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "N1", "C")
		fx.noteRepo.
			On("FindAll", mock.Anything).
			Return([]*domain.Note{n1}, nil).
			Once()

		loadErr := errors.New("db links down")
		fx.linkRepo.
			On("FindAll", mock.Anything).
			Return(([]*domain.Link)(nil), loadErr).
			Once()

		err := fx.svc.LoadGraph(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, loadErr)

		_, err2 := fx.graph.GetNote("n1")
		require.ErrorIs(t, err2, domain.ErrNoteNotFound)
	})

	t.Run("error - no se puede anhadir nota duplicada al grafo.", func(t *testing.T) {
		fx := newGraphSvcMock(t)
		ctx := context.Background()

		existing, _ := domain.NewNote("n1", "existente", "content")
		require.NoError(t, fx.graph.AddNote(existing))

		// Repositorio devuelve la misma nota
		n1dup, _ := domain.NewNote("n1", "duplicada", "content")
		fx.noteRepo.
			On("FindAll", mock.Anything).
			Return([]*domain.Note{n1dup}, nil).
			Once()

		// Link inválido (notas origen/destino no existen)
		fx.linkRepo.
			On("FindAll", mock.Anything).
			Return([]*domain.Link{}, nil).
			Once()

		err := fx.svc.LoadGraph(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNoteAlreadyExists)
	})

	t.Run("No se puede anhadir link sin notas existentes en el grafo.", func(t *testing.T) {
		fx := newGraphSvcMock(t)
		ctx := context.Background()

		// Simulacion: repo de notas vacío
		fx.noteRepo.
			On("FindAll", mock.Anything).
			Return([]*domain.Note{}, nil).
			Once()

		// Link inválido (notas origen/destino no existen)
		l1, _ := domain.NewLink("l1", "n1", "n2", "alias", "title")
		fx.linkRepo.
			On("FindAll", mock.Anything).
			Return([]*domain.Link{l1}, nil).
			Once()

		err := fx.svc.LoadGraph(ctx)
		require.Error(t, err)

		// El dominio debe dispara error
		require.ErrorIs(t, err, domain.ErrOriginNoteNotFound)
	})

}

func TestGraphService_GetConnectedNotes(t *testing.T) {
	t.Run("Éxito - devuelve notas conectadas.", func(t *testing.T) {
		fx := newGraphSvcMock(t)

		// n1 conectado con n2 y n3 (n1->n2, n3->n1)
		n1, _ := domain.NewNote("n1", "N1", "C")
		n2, _ := domain.NewNote("n2", "N2", "C")
		n3, _ := domain.NewNote("n3", "N3", "C")

		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		require.NoError(t, fx.graph.AddNote(n3))

		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}))
		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l2", OriginNoteId: "n3", DestNoteId: "n1"}))

		// Obtener todas notas conectadas a n1
		notes, err := fx.svc.GetConnectedNotes("n1")
		require.NoError(t, err)
		require.Len(t, notes, 2)

		ids := map[string]bool{}
		for _, n := range notes {
			ids[n.ID] = true
		}
		assert.True(t, ids["n2"])
		assert.True(t, ids["n3"])
	})

	t.Run("Error - nota inexistente", func(t *testing.T) {
		fx := newGraphSvcMock(t)

		_, err := fx.svc.GetConnectedNotes("no existe")
		require.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("Error - noteID vacío.", func(t *testing.T) {
		fx := newGraphSvcMock(t)

		_, err := fx.svc.GetConnectedNotes("")
		require.ErrorIs(t, err, domain.ErrInvalidID)
	})
}

func TestGraphService_LinkQueries(t *testing.T) {

	t.Run("Éxito - OutgoingLinks.", func(t *testing.T) {
		fx := newGraphSvcMock(t)

		n1, _ := domain.NewNote("n1", "N1", "C")
		n2, _ := domain.NewNote("n2", "N2", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}))

		links, err := fx.svc.OutgoingLinks("n1")
		require.NoError(t, err)
		require.Len(t, links, 1)
		assert.Equal(t, "l1", links[0].ID)
	})

	t.Run("Éxito - IncomingLinks.", func(t *testing.T) {
		fx := newGraphSvcMock(t)

		n1, _ := domain.NewNote("n1", "N1", "C")
		n2, _ := domain.NewNote("n2", "N2", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}))

		links, err := fx.svc.IncomingLinks("n2")
		require.NoError(t, err)
		require.Len(t, links, 1)
		assert.Equal(t, "l1", links[0].ID)
	})

	t.Run("GetAllLinks delega en grafo", func(t *testing.T) {
		fx := newGraphSvcMock(t)

		n1, _ := domain.NewNote("n1", "N1", "C")
		n2, _ := domain.NewNote("n2", "N2", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}))

		links, err := fx.svc.GetAllLinks("n1")
		require.NoError(t, err)
		require.Len(t, links, 1)
	})

	t.Run("Error - nota inexistente y/o id vacio desde graph", func(t *testing.T) {
		fx := newGraphSvcMock(t)

		_, err := fx.svc.OutgoingLinks("")
		require.ErrorIs(t, err, domain.ErrInvalidID)

		_, err = fx.svc.IncomingLinks("no-existe")
		require.ErrorIs(t, err, domain.ErrNoteNotFound)

		_, err = fx.svc.GetAllLinks("no-existe")
		require.ErrorIs(t, err, domain.ErrNoteNotFound)
	})
}
