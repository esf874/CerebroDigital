package domain

import (
	"testing"
	"time"
)

// Ayuda para punteros (obtener direcciones NoteUpdate)
func ptr[T any](v T) *T { return &v }

func TestNewNote(t *testing.T) {
	// Table-driven test para constructor
	t.Run("Casos de validación - constructor.", func(t *testing.T) {
		tests := []struct {
			name        string
			id          string
			title       string
			expectedErr error
		}{
			{"Válido sin contenido", "id1", "TFG", nil},
			{"Id vacío", "", "Título", ErrInvalidNoteID},
			{"Título vacío", "id1", "", ErrEmptyTitle},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				before := time.Now()
				n, err := NewNote(tt.id, tt.title, "")

				if err != tt.expectedErr {
					t.Fatalf("Esperado error %v, obtenido %v", tt.expectedErr, err)
				}
				if err != nil {
					return
				}

				if n.ID != tt.id {
					t.Errorf("ID incorrecto: esperado %s, obtenido %s", tt.id, n.ID)
				}
				if n.Status != Pending {
					t.Errorf("Estado inicial debe ser Pending, obtenido %s", n.Status)
				}
				if len(n.Tags) != 0 {
					t.Errorf("Tags debe estar vacío al crear, tiene %d elementos", len(n.Tags))
				}
				if n.CreatedAt.Before(before) {
					t.Error("CreatedAt es anterior al inicio del test")
				}
				if !n.CreatedAt.Equal(n.UpdatedAt) {
					t.Error("al crear, CreatedAt y UpdatedAt deben ser iguales")
				}
			})
		}
	})
}

func TestNoteTags(t *testing.T) {
	t.Run("AddTag éxito - actualiza UpdatedAt", func(t *testing.T) {
		n, _ := NewNote("ID2", "T", "C")
		initialUpdate := n.UpdatedAt

		time.Sleep(3 * time.Millisecond)
		err := n.AddTag("Go")

		if err != nil {
			t.Fatalf("Fallo inesperado al añadir tag: %v", err)
		}
		if n.Tags[0] != "Go" {
			t.Errorf("Fallo al añadir tag: %v", n.Tags)
		}
		if !n.UpdatedAt.After(initialUpdate) {
			t.Error("UpdatedAt debería haber avanzado.")
		}
	})

	t.Run("AddTag vacio no modifica la nota.", func(t *testing.T) {
		n, _ := NewNote("ID2", "T", "C")
		before := n.UpdatedAt

		time.Sleep(3 * time.Millisecond)
		err := n.AddTag("")

		if err != ErrEmptyTag {
			t.Errorf("Esperado ErrEmptyTag, obtenido %v", err)
		}

		if !n.UpdatedAt.Equal(before) {
			t.Error("UpdatedAt no debería cambiar sin añadir tag.")
		}
	})

	t.Run("AddTag duplicado, no se anhade, updatedAt no cambia.", func(t *testing.T) {
		n, _ := NewNote("ID2", "T", "C")
		n.AddTag("Go")
		lastUpdate := n.UpdatedAt

		time.Sleep(3 * time.Millisecond)
		err := n.AddTag("Go")

		if err != ErrTagAlreadyExists {
			t.Errorf("Esperado ErrTagAlreadyExists, obtenido %v", err)
		}
		if len(n.Tags) != 1 {
			t.Error("Se ha anhadido erroneamente el tag duplicado.")
		}
		if !n.UpdatedAt.Equal(lastUpdate) {
			t.Error("UpdatedAt NO debería cambiar sin cambios reales")
		}
	})

	t.Run("RemoveTag - éxito con actualización de UpdatedAt", func(t *testing.T) {
		n, _ := NewNote("ID2", "T", "C")
		n.AddTag("Go")
		afterAdd := n.UpdatedAt

		time.Sleep(3 * time.Millisecond)
		err := n.RemoveTag("Go")

		if err != nil {
			t.Fatalf("Fallo inesperado al borrar tag: %v", err)
		}
		if len(n.Tags) != 0 {
			t.Errorf("Tags debe estar vacío tras borrar la nota, tiene %d elementos", len(n.Tags))
		}
		if !n.UpdatedAt.After(afterAdd) {
			t.Error("UpdatedAt debería haber avanzado tras borrar")
		}
	})

	t.Run("RemoveTag - error si tag inexistente", func(t *testing.T) {
		n, _ := NewNote("ID2", "T", "C")

		err := n.RemoveTag("Inexistente")
		if err != ErrTagNotFound {
			t.Errorf("Se esperaba ErrTagNotFound, se obtuvo %v", err)
		}
	})

	t.Run("RemoveTag - al eliminar ultimo elemente, slice vacío.", func(t *testing.T) {
		n, _ := NewNote("ID2", "T", "C")
		_ = n.AddTag("tag")
		_ = n.RemoveTag("tag")

		if n.Tags == nil {
			t.Error("Al eliminar el último tag, Tags no debería ser nil, sino slice vacío.")
		}

		if len(n.Tags) != 0 {
			t.Errorf("Al eliminar el último tag, Tags debería estar vacío, tiene %d elementos", len(n.Tags))
		}

	})

}

func TestNoteApplyUpdate(t *testing.T) {

	t.Run("Update con cambios reales, actualización campos y UpdatedAt", func(t *testing.T) {
		n, _ := NewNote("ID3", "Original", "Contenido")
		initialTime := n.UpdatedAt
		time.Sleep(3 * time.Millisecond)

		update := NoteUpdate{
			Title:    ptr("Nuevo Título"),
			Status:   ptr(Progress),
			Priority: ptr(High),
			Theme:    ptr("TFG"),
		}

		cambiado, err := n.ApplyUpdate(update)
		if err != nil {
			t.Fatalf("Error inesperado en update: %v", err)
		}

		if !cambiado {
			t.Error("Se esperaba 'cambiado' true.")
		}

		// Verificaciones post-update
		if n.Title != "Nuevo Título" || n.Status != Progress || n.Priority != High {
			t.Error("Los campos no se actualizaron correctamente")
		}

		if !n.UpdatedAt.After(initialTime) {
			t.Error("UpdatedAt no cambió tras hacer cambios.")
		}
	})

	t.Run("Update sin cambios reales, no modifica UpdatedAt", func(t *testing.T) {
		n, _ := NewNote("ID3", "Original", "Contenido")
		// Update sin cambios reales
		lastUpdate := n.UpdatedAt
		time.Sleep(3 * time.Millisecond)

		cambiado, err := n.ApplyUpdate(NoteUpdate{
			Title: ptr("Original"),
		})

		if err != nil {
			t.Fatalf("Error inesperado en update sin cambios: %v", err)
		}
		if cambiado {
			t.Error("Se esperaba 'cambiado' false.")
		}
		if !n.UpdatedAt.Equal(lastUpdate) {
			t.Error("UpdatedAt no debería de haber cambiado, no hay cambios.")
		}
	})

	t.Run("Update con struct vacío no modifica nada", func(t *testing.T) {
		n, _ := NewNote("id1", "Título", "Contenido")
		before := n.UpdatedAt
		time.Sleep(3 * time.Millisecond)

		cambiado, err := n.ApplyUpdate(NoteUpdate{})

		if err != nil {
			t.Fatalf("Error inesperado con struct vacío: %v", err)
		}
		if cambiado {
			t.Error("Se esperaba 'cambiado' false con struct vacío.")
		}
		if !n.UpdatedAt.Equal(before) {
			t.Error("UpdatedAt no debería cambiar con struct vacío")
		}
	})

	t.Run("Casos de error de validación", func(t *testing.T) {
		badStatus := NoteStatus("INVALID")
		badPriority := NotePriority("INVALID")

		// Table-driven tests para casos de error
		tests := []struct {
			name        string
			update      NoteUpdate
			expectedErr error
		}{
			{"Título vacío", NoteUpdate{Title: ptr("")}, ErrEmptyTitle},
			{"Status inválido", NoteUpdate{Status: &badStatus}, ErrInvalidStatus},
			{"Priority inválido", NoteUpdate{Priority: &badPriority}, ErrInvalidPriority},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				n, _ := NewNote("id1", "Título", "Contenido")
				before := n.UpdatedAt
				time.Sleep(3 * time.Millisecond)

				cambiado, err := n.ApplyUpdate(tt.update)
				if err != tt.expectedErr {
					t.Errorf("Esperado %v, obtenido %v", tt.expectedErr, err)
				}
				if cambiado {
					t.Error("Se esperaba 'cambiado' false si hay error de validación.")
				}
				if !n.UpdatedAt.Equal(before) {
					t.Error("UpdatedAt no debería cambiar con update inválido")
				}

			})
		}
	})
}
