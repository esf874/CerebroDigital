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

// linkSvcMock es un helper para tests de LinkService, con mocks de repositorios y un grafo real.
type linkSvcMock struct {
	svc      *LinkService
	linkRepo *MockLinkRepository
	noteRepo *MockNoteRepository
	graph    *domain.Graph
}

func newLinkSvcMocks(t *testing.T, idGen IDGenerator) *linkSvcMock {
	t.Helper()

	linkRepo := new(MockLinkRepository)
	noteRepo := new(MockNoteRepository)
	graph := domain.NewGraph()

	svc := NewLinkService(graph, linkRepo, noteRepo, idGen)

	fx := &linkSvcMock{
		svc:      svc,
		linkRepo: linkRepo,
		noteRepo: noteRepo,
		graph:    graph,
	}

	t.Cleanup(func() {
		fx.linkRepo.AssertExpectations(t)
		fx.noteRepo.AssertExpectations(t)
	})

	return fx
}

func TestLinkService_CreateLink(t *testing.T) {

	t.Run("Éxito: sin alías - uso título nota destino.", func(t *testing.T) {
		linkID := "generated-id-link"
		idGen := func() string { return linkID }
		fx := newLinkSvcMocks(t, idGen)
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "origen", "content")
		n2, _ := domain.NewNote("n2", "destino", "content")
		_ = fx.graph.AddNote(n1)
		_ = fx.graph.AddNote(n2)

		fx.linkRepo.On("Save", mock.Anything, mock.MatchedBy(func(link *domain.Link) bool {
			return link.ID == linkID && link.Alias == "destino"
		})).Return(nil).Once()

		link, err := fx.svc.CreateLink(ctx, "n1", "n2", "")

		assert.NoError(t, err)
		assert.Equal(t, link.ID, linkID)

		got, _ := fx.graph.GetLink(linkID)
		assert.Equal(t, "destino", got.Alias)
	})

	t.Run("Éxito: con alías, guarda en repo y anhade al grafo.", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "origen", "content")
		n2, _ := domain.NewNote("n2", "destino", "content")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))

		fx.linkRepo.On("Save", mock.Anything, mock.MatchedBy(func(link *domain.Link) bool {
			return link != nil &&
				link.ID != "" &&
				link.OriginNoteId == "n1" && link.DestNoteId == "n2" &&
				link.Alias == "MiAlias"
		})).Return(nil).Once()

		link, err := fx.svc.CreateLink(ctx, "n1", "n2", "MiAlias")
		require.NoError(t, err)
		require.NotNil(t, link)
		require.Equal(t, "MiAlias", link.Alias)

		got, err := fx.graph.GetLink(link.ID)
		require.NoError(t, err)
		require.Equal(t, "MiAlias", got.Alias)
	})

	t.Run("Error: id nota origen no existe", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		_, err := fx.svc.CreateLink(ctx, "", "n2", "a")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrInvalidID)

		fx.linkRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Error: nota origen no existe", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		nDest, _ := domain.NewNote("nDest", "destino", "content")
		require.NoError(t, fx.graph.AddNote(nDest))

		_, err := fx.svc.CreateLink(ctx, "n1", "n2", "a")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNoteNotFound)

		fx.linkRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Error: nota destino no existe", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		nOrigen, _ := domain.NewNote("nOrigen", "Origen", "C")
		require.NoError(t, fx.graph.AddNote(nOrigen))

		_, err := fx.svc.CreateLink(ctx, "nOrigen", "n2", "a")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNoteNotFound)

		fx.linkRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Error: las notas origen y destino son la misma", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "Origen", "C")
		require.NoError(t, fx.graph.AddNote(n1))

		_, err := fx.svc.CreateLink(ctx, "n1", "n1", "a")
		require.ErrorIs(t, err, domain.ErrSameIds)

		fx.linkRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Fallo simulado en el repositorio", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "Origen", "C")
		n2, _ := domain.NewNote("n2", "DestinoTitle", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))

		saveErr := errors.New("db down")
		fx.linkRepo.
			On("Save", mock.Anything, mock.AnythingOfType("*domain.Link")).
			Return(saveErr).
			Once()

		link, err := fx.svc.CreateLink(ctx, "n1", "n2", "")

		// Servicio devuelve error del repositorio y no crea el link
		require.ErrorIs(t, err, saveErr)
		require.Nil(t, link)

		out, err := fx.graph.OutgoingLinks("n1")
		require.NoError(t, err)
		require.Len(t, out, 0)
	})
}

func TestLinkService_DeleteLink(t *testing.T) {
	t.Run("Éxito: borra en repo y elimina del grafo", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		// Notas origen, destino y link existen en el grafo
		n1, _ := domain.NewNote("n1", "Origen", "C")
		n2, _ := domain.NewNote("n2", "Destino", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))

		l := &domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2", Alias: "a"}
		require.NoError(t, fx.graph.AddLink(l))

		fx.linkRepo.On("Delete", mock.Anything, "l1").Return(nil).Once()

		err := fx.svc.DeleteLink(ctx, "l1")
		require.NoError(t, err)

		// link debe haber sido eliminado del grafo
		_, err = fx.graph.GetLink("l1")
		require.ErrorIs(t, err, domain.ErrLinkNotFound)
	})

	// Si link no está en memoria (por algún fallo), se elimina del repositorio
	t.Run("Éxito: si el link no existe en grafo, se ignora", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		fx.linkRepo.On("Delete", mock.Anything, "l1").Return(nil).Once()

		err := fx.svc.DeleteLink(ctx, "l1")
		require.NoError(t, err)
	})

	// Si falla la BD, el servicio aborta y el link sigue en el grafo/memoria
	t.Run("Error en el repositorio, NO toca grafo", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "Origen", "C")
		n2, _ := domain.NewNote("n2", "Destino", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		l := &domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2", Alias: "a"}
		require.NoError(t, fx.graph.AddLink(l))

		delErr := errors.New("db delete failed")
		fx.linkRepo.
			On("Delete", mock.Anything, "l1").
			Return(delErr).
			Once()

		err := fx.svc.DeleteLink(ctx, "l1")
		require.Error(t, err)
		require.ErrorIs(t, err, delErr)

		_, err2 := fx.graph.GetLink("l1")
		require.NoError(t, err2)
	})
}

func TestLinkService_Queries(t *testing.T) {

	t.Run("Éxito - GetByOrigin", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "Origen", "C")
		n2, _ := domain.NewNote("n2", "Destino", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}))

		// Recupera enlaces salientes
		links, err := fx.svc.GetByOrigin(ctx, "n1")
		require.NoError(t, err)
		require.Len(t, links, 1)
		require.Equal(t, "l1", links[0].ID)
	})

	t.Run("Éxito - GetByDest.", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "Origen", "C")
		n2, _ := domain.NewNote("n2", "Destino", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}))

		// Recupera enlaces entrantes
		links, err := fx.svc.GetByDest(ctx, "n2")
		require.NoError(t, err)
		require.Len(t, links, 1)
		require.Equal(t, "l1", links[0].ID)
	})

	t.Run("Éxito - GetAllForNote.", func(t *testing.T) {
		fx := newLinkSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n1, _ := domain.NewNote("n1", "N1", "C")
		n2, _ := domain.NewNote("n2", "N2", "C")
		require.NoError(t, fx.graph.AddNote(n1))
		require.NoError(t, fx.graph.AddNote(n2))
		require.NoError(t, fx.graph.AddLink(&domain.Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}))

		links, err := fx.svc.GetAllForNote(ctx, "n1")
		require.NoError(t, err)
		require.Len(t, links, 1)
	})

}
