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

// noteSvcMocks es un helper para tests de NoteService, con mocks de repositorios y un grafo real.
type noteSvcMocks struct {
	svc      *NoteService
	noteRepo *MockNoteRepository
	linkRepo *MockLinkRepository
	graph    *domain.Graph
}

func newNoteSvcMocks(t *testing.T, idGen IDGenerator) *noteSvcMocks {
	t.Helper()

	noteRepo := new(MockNoteRepository)
	linkRepo := new(MockLinkRepository)
	graph := domain.NewGraph()
	linkService := NewLinkService(graph, linkRepo, noteRepo, idGen)

	svc := NewNoteService(graph, noteRepo, linkRepo, linkService, idGen)

	fx := &noteSvcMocks{
		svc:      svc,
		noteRepo: noteRepo,
		linkRepo: linkRepo,
		graph:    graph,
	}

	// Se ejecuta automaticamente al fin de cada test, garantiza que cumplan las expectativas mocks (.On())
	t.Cleanup(func() {
		noteRepo.AssertExpectations(t)
		linkRepo.AssertExpectations(t)
	})

	return fx
}

func TestNoteService_CreateNote(t *testing.T) {
	t.Run("Éxito: guarda en repo y añade al grafo", func(t *testing.T) {
		generatedID := "generated"
		idGen := func() string { return generatedID } // Mock funcion generadora
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		fx.noteRepo.
			On("Save", mock.Anything, mock.MatchedBy(func(n *domain.Note) bool {
				return n != nil &&
					n.ID == generatedID && n.Title == "T" &&
					n.Content == "C" && n.Status == domain.Pending
			})).Return(nil).Once()

		n, err := fx.svc.CreateNote(ctx, "T", "C")

		// Validaciones post servicio
		assert.NoError(t, err)
		if assert.NotNil(t, n) {
			assert.Equal(t, generatedID, n.ID)
			assert.Equal(t, domain.Pending, n.Status)
		}

		got, err := fx.graph.GetNote(n.ID)
		assert.NoError(t, err)
		assert.Equal(t, n, got)

		fx.noteRepo.AssertExpectations(t)
	})

	t.Run("Error: valida de dominio (título vacío) -> no llama a repo", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		// Servicio intenta crear nota, dominio debe lanzar error
		n, err := fx.svc.CreateNote(ctx, "", "C")
		assert.ErrorIs(t, err, domain.ErrEmptyTitle)
		assert.Nil(t, n)

		fx.noteRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		fx.noteRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
		fx.linkRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Fallo en BD -> no se modifica grafo", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()
		saveErr := errors.New("DB down")

		// Simulación de error
		fx.noteRepo.
			On("Save", mock.Anything, mock.AnythingOfType("*domain.Note")).Return(saveErr).Once()
		note, err := fx.svc.CreateNote(ctx, "T", "C")

		// Si falla persistencia, servicio devuelve error sin modificar grafo
		assert.ErrorIs(t, err, saveErr)
		assert.Nil(t, note)
		assert.Empty(t, fx.graph.GetAllNotes())
	})

	t.Run("Error - id duplicado, rollback en repo.", func(t *testing.T) {
		generatedID := "generated"
		idGen := func() string { return generatedID }

		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		// Simular id duplicado en el grafo
		note, _ := domain.NewNote(generatedID, "existente", "content")
		_ = fx.graph.AddNote(note)

		fx.noteRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

		// rollback - si tras guardar en BD no se anhade al grafo, el servicio borra de BD
		fx.noteRepo.On("Delete", mock.Anything, generatedID).Return(nil).Once()

		note2, err := fx.svc.CreateNote(ctx, "nuevo titulo", "new content")

		assert.Error(t, err)
		assert.Nil(t, note2)

		// Verifica que grafo tiene nota original
		original, _ := fx.graph.GetNote(generatedID)
		assert.Equal(t, note, original)
	})
}

func TestNoteService_GetNote(t *testing.T) {
	t.Run("Éxito: devuelve nota existente en grafo", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		// Insertar nota
		n, err := domain.NewNote("n1", "T", "C")
		require.NoError(t, err)
		err = fx.graph.AddNote(n)
		require.NoError(t, err)

		// Servicio recupera nota de memoria
		got, err := fx.svc.GetNote(ctx, "n1")
		assert.NoError(t, err)
		if assert.NotNil(t, got) {
			assert.Equal(t, "n1", got.ID)
		}
	})

	t.Run("Error: id vacío.", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		got, err := fx.svc.GetNote(ctx, "")
		assert.ErrorIs(t, err, domain.ErrInvalidID)
		assert.Nil(t, got)
	})

	t.Run("Error: nota inexistente.", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		got, err := fx.svc.GetNote(ctx, "no existe")
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
		assert.Nil(t, got)
	})
}

func TestNoteService_UpdateNote(t *testing.T) {
	t.Run("Éxito: aplica update y persiste.", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n, err := domain.NewNote("n1", "Old", "C")
		require.NoError(t, err)
		err = fx.graph.AddNote(n)
		require.NoError(t, err)

		// Valida que la logica de applyUpdate es correcta
		fx.noteRepo.On("Save", mock.Anything, mock.MatchedBy(func(n *domain.Note) bool {
			return n != nil && n.ID == "n1" && n.Title == "New"
		})).Return(nil).Once()

		update := domain.NoteUpdate{Title: ptr("New")}
		err = fx.svc.UpdateNote(ctx, "n1", update)

		assert.NoError(t, err)

		// En memoria se refleja el cambio
		got, err := fx.graph.GetNote("n1")
		assert.NoError(t, err)
		if assert.NotNil(t, got) {
			assert.Equal(t, "New", got.Title)
		}
	})

	t.Run("Error: título vacío -> no se persiste.", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n, err := domain.NewNote("n1", "Old", "C")
		require.NoError(t, err)
		err = fx.graph.AddNote(n)
		require.NoError(t, err)

		update := domain.NoteUpdate{Title: ptr("")}
		err = fx.svc.UpdateNote(ctx, "n1", update)
		assert.ErrorIs(t, err, domain.ErrEmptyTitle)

		fx.noteRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)

		// La nota no debe cambiar
		got, err := fx.graph.GetNote("n1")
		assert.NoError(t, err)
		if assert.NotNil(t, got) {
			assert.Equal(t, "Old", got.Title)
		}
	})

	t.Run("Error: fallo del repositorio", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n, _ := domain.NewNote("n1", "Old", "C")
		_ = fx.graph.AddNote(n)

		saveErr := errors.New("write failed")
		fx.noteRepo.
			On("Save", mock.Anything, mock.Anything).Return(saveErr).Once()

		update := domain.NoteUpdate{Title: ptr("New")}
		err := fx.svc.UpdateNote(ctx, "n1", update)

		assert.ErrorIs(t, err, saveErr)

		// Rollback, la nota no debe cambiar en el grafo
		got, _ := fx.graph.GetNote("n1")
		assert.Equal(t, "Old", got.Title)
	})

	t.Run("Update sin cambios reales: no debe persistir en repo", func(t *testing.T) {
		fx := newNoteSvcMocks(t, func() string { return "generated" })
		ctx := context.Background()

		n, _ := domain.NewNote("n1", "Igual", "C")
		_ = fx.graph.AddNote(n)

		update := domain.NoteUpdate{Title: ptr("Igual")}

		err := fx.svc.UpdateNote(ctx, "n1", update)
		assert.NoError(t, err)

		// Si dominio detecta que NO hay cambios, no se llama a infraestructura
		fx.noteRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})
}

func TestNoteService_DeleteNote(t *testing.T) {

	idGen := func() string { return "generated" }
	t.Run("Éxito: borra links automáticamente, borra nota en repo y en grafo", func(t *testing.T) {
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		n, _ := domain.NewNote("n1", "T", "C")
		_ = fx.graph.AddNote(n)

		// Orden: primero borra links asociados, luego persiste y final memoria
		fx.linkRepo.On("DeleteAllByNoteID", mock.Anything, "n1").Return(nil).Once()
		fx.noteRepo.On("Delete", mock.Anything, "n1").Return(nil).Once()

		err := fx.svc.DeleteNote(ctx, "n1")
		assert.NoError(t, err)

		// Verificar que la nota se elimino del grafo/memoria
		_, err = fx.graph.GetNote("n1")
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})

	t.Run("Error: id vacío => ErrInvalidID", func(t *testing.T) {
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		err := fx.svc.DeleteNote(ctx, "")
		assert.ErrorIs(t, err, domain.ErrInvalidID)

		fx.noteRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
		fx.linkRepo.AssertNotCalled(t, "FindByNoteID", mock.Anything, mock.Anything)
	})

	t.Run("Error: linkRepo.FindByNoteID falla", func(t *testing.T) {
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		delErr := errors.New("delete failed")
		fx.linkRepo.
			On("DeleteAllByNoteID", mock.Anything, "n1").
			Return(delErr).
			Once()

		err := fx.svc.DeleteNote(ctx, "n1")
		assert.ErrorContains(t, err, delErr.Error())
		fx.noteRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("Error: borrar un link falla => aborta", func(t *testing.T) {
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		delErr := errors.New("Delete failed")

		fx.linkRepo.
			On("DeleteAllByNoteID", mock.Anything, "n1").
			Return(delErr).
			Once()

		err := fx.svc.DeleteNote(ctx, "n1")
		assert.ErrorContains(t, err, delErr.Error())

		// Si fallan los links NO debe intentar borrar la nota
		fx.noteRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("Error: noteRepo.Delete falla => la nota sigue en el grafo", func(t *testing.T) {
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		n, err := domain.NewNote("n1", "T", "C")
		err = fx.graph.AddNote(n)

		fx.linkRepo.
			On("DeleteAllByNoteID", mock.Anything, "n1").
			Return(nil).
			Once()

		// Simular error en el repositorio
		delErr := errors.New("db delete failed")
		fx.noteRepo.
			On("Delete", mock.Anything, "n1").
			Return(delErr).
			Once()

		err = fx.svc.DeleteNote(ctx, "n1")
		assert.ErrorIs(t, err, delErr)

		// La nota sigue en el grafo (no se llamó a RemoveNote)
		_, err = fx.graph.GetNote("n1")
		assert.NoError(t, err)
	})
}

func TestNoteService_Search(t *testing.T) {

	idGen := func() string { return "generated" }

	t.Run("Éxito - SearchByContent.", func(t *testing.T) {
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		// Simula la recuperacion de notas de consulta
		expected := []*domain.Note{{ID: "n1"}, {ID: "n2"}}

		fx.noteRepo.
			On("SearchByContent", mock.Anything, "hello").Return(expected, nil).Once()

		got, err := fx.svc.Search(ctx, "hello")
		assert.NoError(t, err)
		if assert.NotNil(t, got) {
			assert.Len(t, got, 2) 
			assert.Equal(t, "n1", got[0].ID)
		}
	})

	t.Run("Éxito - fidByTag", func(t *testing.T) {
		fx := newNoteSvcMocks(t, idGen)
		ctx := context.Background()

		expected := []*domain.Note{{ID: "n1"}}

		fx.noteRepo.
			On("FindByTag", mock.Anything, "go").
			Return(expected, nil).
			Once()

		got, err := fx.svc.SearchByTag(ctx, "go")
		assert.NoError(t, err)
		if assert.NotNil(t, got) {
			assert.Len(t, got, 1)
		}
	})
}

// Helper para punteros (NoteUpdate)
func ptr[T any](v T) *T { return &v }
