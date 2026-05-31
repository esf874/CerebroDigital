package domain

import (
	"fmt"
	"sync"
	"testing"
)

// Helper para construir un grafo con notas dadas.
func newGraphWithNotes(t *testing.T, ids ...string) *Graph {
	t.Helper()
	g := NewGraph()
	for _, id := range ids {
		n, _ := NewNote(id, "Título "+id, "Contenido")
		if err := g.AddNote(n); err != nil {
			t.Fatalf("Setup: error al anhadir nota %s: %v", id, err)
		}
	}
	return g
}

// Helper para anhadir un link.
func addLink(t *testing.T, g *Graph, id, origin, dest string) {
	t.Helper()
	if err := g.AddLink(&Link{ID: id, OriginNoteId: origin, DestNoteId: dest}); err != nil {
		t.Fatalf("Setup: error al link %s: %v", id, err)
	}
}

func TestNewGraph(t *testing.T) {
	g := NewGraph()

	if g == nil {
		t.Fatal("NewGraph no deberia devolver nil.")
	}

	// Inicialización correcto de los mapas.
	if g.Notes == nil {
		t.Fatal("Mapa de notes no debería ser nil.")
	}
	if g.Links == nil {
		t.Fatal("Mapa de links no debería ser nil.")
	}

	if (len(g.Notes) != 0) || (len(g.Links) != 0) {
		t.Error("Los mapas deberían estar vacíos al crear el grafo.")
	}
}

func TestGraph_AddNote(t *testing.T) {
	t.Run("Caso de éxito - addNote", func(t *testing.T) {
		g := NewGraph()
		n1, _ := NewNote("n1", "Nota 1", "Contenido")

		err := g.AddNote(n1)
		if err != nil {
			t.Fatalf("No se esperaba error al añadir nota: %v", err)
		}
		if len(g.Notes) != 1 {
			t.Errorf("El grafo debería tener 1 nota, hay %d", len(g.Notes))
		}
	})

	t.Run("Error - id vacío.", func(t *testing.T) {
		noteEmptyId, _ := NewNote("temp", "Título", "C")
		noteEmptyId.ID = ""

		tests := []struct {
			name        string
			note        *Note
			expectedErr error
		}{
			{"Nota nil", nil, ErrNilNote},
			{"Nota id vacío", noteEmptyId, ErrInvalidID},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				g := NewGraph()
				err := g.AddNote(tt.note)
				if err != tt.expectedErr {
					t.Errorf("Esperado error %v, obtenido %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("Error - nota duplicada.", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1")
		n1, _ := NewNote("n1", "Duplicada", "")

		err := g.AddNote(n1)
		if err != ErrNoteAlreadyExists {
			t.Errorf("Esperado ErrNoteAlreadyExists, obtenido %v", err)
		}
	})
}

func TestGraph_GetNote(t *testing.T) {
	t.Run("Caso de éxito", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1")
		found, err := g.GetNote("n1")

		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if found == nil || found.ID != "n1" {
			t.Error("Se esperaba recuperar ntoa n1.")
		}
	})

	t.Run("Casos de error GetNote", func(t *testing.T) {
		tests := []struct {
			name        string
			id          string
			expectedErr error
		}{
			{"ID vacío", "", ErrInvalidID},
			{"Nota inexistente", "no existe", ErrNoteNotFound},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				g := NewGraph()
				_, err := g.GetNote(tt.id)
				if err != tt.expectedErr {
					t.Errorf("Esperado error %v, obtenido %v", tt.expectedErr, err)
				}
			})
		}
	})
}

func TestGraph_GetAllNotes(t *testing.T) {
	t.Run("Grafo vacío devuelve slice vacío", func(t *testing.T) {
		g := NewGraph()
		notes := g.GetAllNotes()
		if notes == nil {
			t.Error("GetAllNotes no debería devolver nil, debería devolver slice vacío.")
		}
		if len(notes) != 0 {
			t.Errorf("Esperado 0 notas, hay %d", len(notes))
		}
	})

	t.Run("Grafo con notas devuelve todas las notas", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1", "n2", "n3")
		notes := g.GetAllNotes()
		if len(notes) != 3 {
			t.Errorf("Esperado 3 notas, hay %d", len(notes))
		}
	})
}

func TestGraph_RemoveNote(t *testing.T) {
	t.Run("Caso de éxito remove sin links asociados", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1")
		err := g.RemoveNote("n1")

		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if len(g.Notes) != 0 {
			t.Errorf("Tras borrar, notes debería tener 0 notas.")
		}
	})

	t.Run("Caso de éxito remove con links asociados.", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1", "n2", "n3")
		addLink(t, g, "l1", "n1", "n2")
		addLink(t, g, "l2", "n2", "n1")
		addLink(t, g, "l3", "n2", "n3")

		err := g.RemoveNote("n1")
		if err != nil {
			t.Fatalf("Error al eliminar nota: %v", err)
		}

		// Verificar que la nota fue eliminada
		if _, err := g.GetNote("n1"); err != ErrNoteNotFound {
			t.Error("la nota n1 debería haber sido eliminada del grafo")
		}

		// l1 y l2 deben haber sido borrados
		if _, err := g.GetLink("l1"); err != ErrLinkNotFound {
			t.Error("El link debería haber sido eliminado.")
		}
		if _, err := g.GetLink("l2"); err != ErrLinkNotFound {
			t.Error("El link debería haber sido eliminado.")
		}
		// l3 no debe haber sido borrado
		if _, err := g.GetLink("l3"); err != nil {
			t.Errorf("l3 debería existir; err=%v", err)
		}
	})

	t.Run("Casos de error RemoveNote", func(t *testing.T) {
		tests := []struct {
			name        string
			id          string
			expectedErr error
		}{
			{"ID vacío", "", ErrInvalidID},
			{"Nota inexistente", "no existe", ErrNoteNotFound},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				g := NewGraph()
				err := g.RemoveNote(tt.id)
				if err != tt.expectedErr {
					t.Errorf("Esperado error %v, obtenido %v", tt.expectedErr, err)
				}
			})
		}
	})
}

func TestGraph_AddLink(t *testing.T) {
	t.Run("Caso de éxito", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1", "n2")
		l := &Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"}

		err := g.AddLink(l)
		if err != nil {
			t.Fatalf("Error inesperado al añadir link válido: %v", err)
		}
		if len(g.Links) != 1 {
			t.Errorf("Debería de haber un link, hay %d .", len(g.Links))
		}
	})

	t.Run("Casos de error AddLink", func(t *testing.T) {
		tests := []struct {
			name        string
			link        *Link
			expectedErr error
		}{
			{"link nil", nil, ErrNilLink},
			{"id link vacío", &Link{ID: "", OriginNoteId: "n1", DestNoteId: "n2"}, ErrInvalidLinkID},
			{"origin vacío", &Link{ID: "l1", OriginNoteId: "", DestNoteId: "n2"}, ErrInvalidID},
			{"dest vacío", &Link{ID: "l1", OriginNoteId: "n1", DestNoteId: ""}, ErrInvalidID},
			{"origin == dest", &Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n1"}, ErrSameIds},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				g := newGraphWithNotes(t, "n1", "n2")
				err := g.AddLink(tt.link)
				if err != tt.expectedErr {
					t.Errorf("esperado %v, obtenido %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("Nota origen no existe en el grafo.", func(t *testing.T) {
		g := newGraphWithNotes(t, "n2")
		l := &Link{ID: "l1", OriginNoteId: "fake", DestNoteId: "n2"}

		err := g.AddLink(l)
		if err != ErrOriginNoteNotFound {
			t.Errorf("Esperado ErrOriginNoteNotFound, obtenido %v", err)
		}
	})

	t.Run("Nota destino no existe en el grafo.", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1")
		l := &Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "fake"}

		err := g.AddLink(l)
		if err != ErrDestNoteNotFound {
			t.Errorf("Esperado ErrDestNoteNotFound, obtenido %v", err)
		}
	})

	t.Run("Links duplicados", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1", "n2")
		addLink(t, g, "l1", "n1", "n2")

		err := g.AddLink(&Link{ID: "l1", OriginNoteId: "n1", DestNoteId: "n2"})

		if err != ErrLinkAlreadyExists {
			t.Errorf("Esperado ErrLinkAlreadyExists, obtenido %v.", err)
		}
	})
}

func TestGraph_GetLink(t *testing.T) {
	t.Run("Caso de éxito", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1", "n2")
		addLink(t, g, "l1", "n1", "n2")

		l, err := g.GetLink("l1")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if l == nil || l.ID != "l1" {
			t.Error("Se esperaba recuperar el link l1.")
		}
	})

	t.Run("Casos de error GetLink", func(t *testing.T) {
		tests := []struct {
			name        string
			id          string
			expectedErr error
		}{
			{"ID vacío", "", ErrInvalidLinkID},
			{"Link inexistente", "no existe", ErrLinkNotFound},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				g := NewGraph()
				_, err := g.GetLink(tt.id)
				if err != tt.expectedErr {
					t.Errorf("Esperado error %v, obtenido %v", tt.expectedErr, err)
				}
			})
		}
	})
}

func TestGraph_RemoveLink(t *testing.T) {
	t.Run("Caso éxito.", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1", "n2")
		addLink(t, g, "l1", "n1", "n2")

		err := g.RemoveLink("l1")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if _, err := g.GetLink("l1"); err != ErrLinkNotFound {
			t.Error("Tras borrar, el link debería haber sido eliminado.")
		}
	})

	t.Run("Casos de error RemoveLink", func(t *testing.T) {
		tests := []struct {
			name        string
			id          string
			expectedErr error
		}{
			{"id vacío", "", ErrInvalidLinkID},
			{"link inexistente", "no-existe", ErrLinkNotFound},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				g := NewGraph()
				err := g.RemoveLink(tt.id)
				if err != tt.expectedErr {
					t.Errorf("Esperado error %v, obtenido %v", tt.expectedErr, err)
				}
			})
		}
	})
}

func TestGraph_LinkQueries(t *testing.T) {
	/** n1->n2, n1->n3, n2->n1
	n1: 2 salientes, 1 entrante, 3 totales
	n2: 1 saliente, 1 entrante, 2 totales
	n3: 0 salientes, 1 entrante, 1 total */

	setup := func(t *testing.T) *Graph {
		t.Helper()
		g := newGraphWithNotes(t, "n1", "n2", "n3")
		addLink(t, g, "l1", "n1", "n2")
		addLink(t, g, "l2", "n1", "n3")
		addLink(t, g, "l3", "n2", "n1")
		return g
	}

	t.Run("OutgoingLinks", func(t *testing.T) {
		g := setup(t)

		// Salientes - n1 debe tener 2 aristas
		out, err := g.OutgoingLinks("n1")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if len(out) != 2 {
			t.Errorf("n1 debería tener 2 enlaces salientes, tiene %d", len(out))
		}
	})

	t.Run("IncomingLinks", func(t *testing.T) {
		g := setup(t)

		// Entrantes - n1 debe tener 1
		in, err := g.IncomingLinks("n1")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if len(in) != 1 {
			t.Errorf("n1 debería tener 1 enlace de entrada, tiene %d", len(in))
		}
	})

	t.Run("AllLinksForNote", func(t *testing.T) {
		g := setup(t)

		// n1 debe tener 3 en total
		all, err := g.AllLinksForNote("n1")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if len(all) != 3 {
			t.Errorf("n1 debería tener 3 enlaces totales, tiene %d", len(all))
		}
	})

	t.Run("AllLinksForNote - nota sin links, slice vacío.", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1")

		all, err := g.AllLinksForNote("n1")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}
		if all == nil {
			t.Error("Deberia devolver slice vacío, no nil.")
		}
		if len(all) != 0 {
			t.Errorf("n1 no debería tener links, tiene %d .", len(all))
		}
	})

	t.Run("Nota inexistente.", func(t *testing.T) {
		g := NewGraph()

		if _, err := g.OutgoingLinks("inexistente"); err != ErrNoteNotFound {
			t.Errorf("Enlaces salientes: esperado ErrNoteNotFound, obtenido %v", err)
		}

		if _, err := g.IncomingLinks("inexistente"); err != ErrNoteNotFound {
			t.Errorf("Enlaces entrantes: esperado ErrNoteNotFound, obtenido %v", err)
		}

		if _, err := g.AllLinksForNote("inexistente"); err != ErrNoteNotFound {
			t.Errorf("Todos los enlaces: esperado ErrNoteNotFound, obtenido %v", err)
		}
	})

	t.Run("ID vacío.", func(t *testing.T) {
		g := NewGraph()

		if _, err := g.OutgoingLinks(""); err != ErrInvalidID {
			t.Errorf("Enlaces salientes: esperado ErrInvalidID, obtenido %v", err)
		}

		if _, err := g.IncomingLinks(""); err != ErrInvalidID {
			t.Errorf("Enlaces entrantes: esperado ErrInvalidID, obtenido %v", err)
		}

		if _, err := g.AllLinksForNote(""); err != ErrInvalidID {
			t.Errorf("Todos los enlaces: esperado ErrInvalidID, obtenido %v", err)
		}
	})
}

func TestGraph_Concurrency(t *testing.T) {

	// Muchas goroutines pueden escribir simultaneamente
	t.Run("Escrituras concurrentes.", func(t *testing.T) {
		g := NewGraph()
		var waitGroup sync.WaitGroup
		workers := 15

		for i := 0; i < workers; i++ {
			waitGroup.Add(1)
			go func(i int) {
				defer waitGroup.Done() // Indica que goroutine terminó
				noteID := fmt.Sprintf("n%d", i)
				n, _ := NewNote(noteID, "Título", "Contenido")

				g.AddNote(n) // Adquiere el mutex
			}(i) 
		}
		waitGroup.Wait() 

		// El numero de notas debe ser igual al de workers
		if len(g.Notes) != workers {
			t.Errorf("Esperado %d notas, hay %d", workers, len(g.Notes))
		}
	})

	// Lecturas concurrentes mientras se escribe
	// rlock permite muchos lectores a la vez pero si una goroutine solicita escritura,
	// bloquea hasta que todos los lectores liberen el lock
	t.Run("Lecturas y escrituras concurrentes.", func(t *testing.T) {
		g := newGraphWithNotes(t, "n1", "n2")
		addLink(t, g, "l1", "n1", "n2")
		var waitGroup sync.WaitGroup

		for i := 0; i < 20; i++ {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				g.GetNote("n1")       
				g.OutgoingLinks("n1") 
				g.GetAllNotes()     
			}()
		}

		// Mientras hay lectores activos, lanza una goroutine escritura
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			n, _ := NewNote("n3", "nueva", "")
			_ = g.AddNote(n)
		}()
		waitGroup.Wait()
	})
}
