package domain

import "sync"

// Graph representa el grafo en memoria, gestiona notas, enlaces y un índice de tags para búsquedas eficientes.
type Graph struct {
	// Mutex que garantiza el acceso concurrente al grafo.
	mutex sync.RWMutex
	Notes map[string]*Note
	Links map[string]*Link

	// Búsqueda por tags:key=tag, valor=id notas que lo contienen
	TagsIndex map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		Notes:     make(map[string]*Note),
		Links:     make(map[string]*Link),
		TagsIndex: make(map[string][]string),
	}
}

// Métodos gestión lógica interna del grafo
func (g *Graph) AddNote(note *Note) error {
	if note == nil {
		return ErrNilNote
	}
	if note.ID == "" {
		return ErrInvalidID
	}

	g.mutex.Lock()    
	defer g.mutex.Unlock() 

	// Nota ya existente
	if _, exists := g.Notes[note.ID]; exists {
		return ErrNoteAlreadyExists
	} 

	g.Notes[note.ID] = note

	for _, tag := range note.Tags {
		g.TagsIndex[tag] = append(g.TagsIndex[tag], note.ID)
	}
	return nil
}

func (g *Graph) GetNote(id string) (*Note, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	g.mutex.RLock() 
	defer g.mutex.RUnlock()

	note, exists := g.Notes[id]
	if !exists {
		return nil, ErrNoteNotFound
	}
	return note, nil
}

func (g *Graph) RemoveNote(id string) error {
	if id == "" {
		return ErrInvalidID
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	note, exists := g.Notes[id]
	if !exists {
		return ErrNoteNotFound
	}

	// Actualizar índice de tags: elimina id de nota de cada tag
	for _, tag := range note.Tags {
		ids := g.TagsIndex[tag]
		for i, noteID := range ids {
			if noteID == id {
				g.TagsIndex[tag] = append(ids[:i], ids[i+1:]...)
				break
			}
		}
	}

	// Eliminar links asociados
	for linkID, link := range g.Links {
		if link.OriginNoteId == id || link.DestNoteId == id {
			delete(g.Links, linkID)
		}
	}

	delete(g.Notes, id)
	return nil
}

func (g *Graph) RemoveLinksByOrigin(originID string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for id, link := range g.Links {
		if link.OriginNoteId == originID {
			delete(g.Links, id)
		}
	}
}

func (g *Graph) GetNotesByTag(tag string) []string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	// Retorna copia del slice para evitar modificacion
	ids, exists := g.TagsIndex[tag]
	if !exists {
		return []string{}
	}

	result := make([]string, len(ids))
	copy(result, ids)
	return result
}

func (g *Graph) GetAllNotes() []*Note {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	notes := make([]*Note, 0, len(g.Notes))
	for _, note := range g.Notes {
		notes = append(notes, note)
	}
	return notes
}

func (g *Graph) AddLink(link *Link) error {
	if link == nil {
		return ErrNilLink
	}

	if link.ID == "" {
		return ErrInvalidLinkID
	}

	if link.OriginNoteId == "" || link.DestNoteId == "" {
		return ErrInvalidID
	}

	if link.OriginNoteId == link.DestNoteId {
		return ErrSameIds
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Comprobar existencia de notas origen y destino
	if _, exists := g.Notes[link.OriginNoteId]; !exists {
		return ErrOriginNoteNotFound
	}
	if _, exists := g.Notes[link.DestNoteId]; !exists {
		return ErrDestNoteNotFound
	}

	// Comprobar duplicados
	if _, exists := g.Links[link.ID]; exists {
		return ErrLinkAlreadyExists
	}

	g.Links[link.ID] = link
	return nil
}

// Gestión links
func (g *Graph) GetLink(id string) (*Link, error) {
	if id == "" {
		return nil, ErrInvalidLinkID
	}

	g.mutex.RLock()
	defer g.mutex.RUnlock()

	link, exists := g.Links[id]
	if !exists {
		return nil, ErrLinkNotFound
	}
	return link, nil
}

func (g *Graph) OutgoingLinks(noteID string) ([]*Link, error) {
	if noteID == "" {
		return nil, ErrInvalidID
	}

	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if _, exists := g.Notes[noteID]; !exists {
		return nil, ErrNoteNotFound
	}

	result := []*Link{}
	for _, link := range g.Links {
		if link.OriginNoteId == noteID {
			result = append(result, link)
		}
	}
	return result, nil
}

func (g *Graph) IncomingLinks(noteID string) ([]*Link, error) {
	if noteID == "" {
		return nil, ErrInvalidID
	}
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if _, exists := g.Notes[noteID]; !exists {
		return nil, ErrNoteNotFound
	}

	result := []*Link{}
	for _, link := range g.Links {
		if link.DestNoteId == noteID {
			result = append(result, link)
		}
	}
	return result, nil
}

func (g *Graph) RemoveLink(id string) error {
	if id == "" {
		return ErrInvalidLinkID
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.Links[id]; !exists {
		return ErrLinkNotFound
	}
	delete(g.Links, id)
	return nil
}

func (g *Graph) AllLinksForNote(noteID string) ([]*Link, error) {
	if noteID == "" {
		return nil, ErrInvalidID
	}

	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if _, exists := g.Notes[noteID]; !exists {
		return nil, ErrNoteNotFound
	}

	result := []*Link{}
	for _, link := range g.Links {
		if link.OriginNoteId == noteID || link.DestNoteId == noteID {
			result = append(result, link)
		}
	}
	return result, nil
}
