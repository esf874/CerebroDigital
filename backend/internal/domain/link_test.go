package domain

import "testing"
 
func TestNewLink(t *testing.T) {
	t.Run("Creacion link con alias vacío", func(t *testing.T) {

		link, err := NewLink("l1", "n1", "n2", "", "Título nota destino")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}

		if link.Alias != "Título nota destino" {
			t.Errorf("Alias incorrecto, esperado 'Título nota destino', obtenido %q", link.Alias)
		}
	})

	t.Run("Creacion link con alías", func(t *testing.T) {
		l, err := NewLink("l1", "n1", "n2", "depende de", "título ignorado")
		if err != nil {
			t.Fatalf("No se esperaba error: %v", err)
		}

		if l.Alias != "depende de" {
			t.Errorf("Alias incorrecto, esperado %q, obtenido %q", "depende de", l.Alias)
		}
	})

	t.Run("Validaciones de errores con id.", func(t *testing.T) {
		tests := []struct {
			name     string
			id, o, d string
			want     error
		}{
			{"ID vacío", "", "n1", "n2", ErrInvalidID},
			{"Origen vacío", "l1", "", "n2", ErrInvalidID},
			{"Destino vacío", "l1", "n1", "", ErrInvalidID},
			{"Mismo ID", "l1", "n1", "n1", ErrSameIds},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := NewLink(tt.id, tt.o, tt.d, "alias", "titulo")
				if err != tt.want {
					t.Errorf("Esperado %v, obtenido %v", tt.want, err)
				}
			})
		}
	})
}
