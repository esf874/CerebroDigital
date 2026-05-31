package application

import (
	"context"
	"fmt"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
)

// GraphService coordina las operaciones sobre el grafo, incluyendo la carga inicial desde 
// los repositorios y las consultas de navegación.
type GraphService struct {
	graph    *domain.Graph
	noteRepo domain.InterfaceNoteRepository
	linkRepo domain.InterfaceLinkRepository
}

// NewGraphService crea el servicio de grafo e inyecta sus dependencias.
func NewGraphService(
	graph *domain.Graph, 
	noteRepo domain.InterfaceNoteRepository,
	linkRepo domain.InterfaceLinkRepository,
) *GraphService {
	return &GraphService{
		graph:    graph,
		noteRepo: noteRepo,
		linkRepo: linkRepo,
	}
}

// LoadGraph reconstruye el grafo en memoria a partir de los repositorios
// durante el arranque de la aplicación.
func (g *GraphService) LoadGraph(ctx context.Context) error {

	notes, err := g.noteRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("Failed to load notes: %w", err)
	}

	links, err := g.linkRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("Failed to load links: %w", err)
	}

	for _, note := range notes {
		if err := g.graph.AddNote(note); err != nil {
			return fmt.Errorf("Failed to add note %s: %w", note.ID, err)
		}
	}

	for _, link := range links {
		if err := g.graph.AddLink(link); err != nil {
			return fmt.Errorf("Failed to add link %s: %w", link.ID, err)
		}

	}
	return nil
}

// GetConnectedNotes devuelve las notas conectadas a una nota dada, ya sea por enlaces explícitos o por tags compartidos.
func (s *GraphService) GetConnectedNotes(noteID string) ([]*domain.Note, error) {

	links, err := s.graph.AllLinksForNote(noteID)
	if err != nil {
		return nil, err
	}

	// Uso de mapa auxiliar para evitar duplicados (enlaces bidireccionales)
	noteIDs := make(map[string]bool)
	for _, link := range links {
		if link.OriginNoteId == noteID {
			noteIDs[link.DestNoteId] = true
		}
		if link.DestNoteId == noteID {
			noteIDs[link.OriginNoteId] = true
		}
	}

	notes := make([]*domain.Note, 0, len(noteIDs))
	for id := range noteIDs {
		note, err := s.graph.GetNote(id)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (g *GraphService) OutgoingLinks(noteID string) ([]*domain.Link, error) {
	return g.graph.OutgoingLinks(noteID)
}

func (g *GraphService) IncomingLinks(noteID string) ([]*domain.Link, error) {
	return g.graph.IncomingLinks(noteID)
}

func (g *GraphService) GetAllLinks(noteID string) ([]*domain.Link, error) {
	return g.graph.AllLinksForNote(noteID)
}

// GraphNodeDTO representa un nodo del grafo enviado al frontend.
type GraphNodeDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	IsCurrent bool   `json:"isCurrent"`
	Depth     int    `json:"depth"`     
}

// GraphLinkDTO representa una arista del grafo enviada al frontend.
type GraphLinkDTO struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// GraphResponse representa la respuesta de la consulta de subgrafo, con nodos y aristas.
type GraphResponse struct {
	Nodes []GraphNodeDTO `json:"nodes"`
	Edges []GraphLinkDTO `json:"edges"`
}

// Estructura interna auxiliar para evitar duplicados 
type graphEdge struct {
	source string
	target string
}

// bfsNode representa un nodo en la búsqueda en anchura, con ID y profundidad relativa al nodo raíz. 
type bfsNode struct {
	id    string
	depth int
}

// bfsResult representa el resultado de la búsqueda en anchura.
type bfsResult struct {
	Nodes []bfsNode
	Edges []graphEdge
}

// normalizeEdgeKey genera una clave única para una arista entre dos nodos, independientemente del sentido de la arista.
func normalizeEdgeKey(id1, id2 string) string {
	if id1 < id2 {
		return id1 + "-" + id2
	}
	return id2 + "-" + id1
}

// bfs realiza una búsqueda en anchura desde la nota origen.
// Explora tanto enlaces explícitos como relaciones semánticas por tags, hasta una profundidad y cantidad de nodos limitados.
func (g *GraphService) bfs(noteID string, maxDepth int, maxNodes int) (*bfsResult, error) {
	_, err := g.graph.GetNote(noteID)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	queue := []bfsNode{{id: noteID, depth: 0}}
	visited[noteID] = true

	var nodes []bfsNode
	nodes = append(nodes, bfsNode{id: noteID, depth: 0})

	// Mapa para evitar aristas duplicados
	edgesSet := make(map[string]graphEdge)

	for len(queue) > 0 && len(nodes) < maxNodes {
		current := queue[0]
		queue = queue[1:]

		if current.depth >= maxDepth {
			continue
		}

		currentNote, _ := g.graph.GetNote(current.id)

		// Exploracion relaciones explícitas
		links, _ := g.graph.AllLinksForNote(current.id)
		for _, link := range links {
			neighborID := link.DestNoteId
			if link.DestNoteId == current.id {
				neighborID = link.OriginNoteId
			}

			//  Evita duplicados por link bidireccional
			edgeKey := normalizeEdgeKey(current.id, neighborID)
			edgesSet[edgeKey] = graphEdge{source: current.id, target: neighborID}

			if !visited[neighborID] && len(nodes) < maxNodes {
				visited[neighborID] = true
				newNode := bfsNode{id: neighborID, depth: current.depth + 1}
				queue = append(queue, newNode)
				nodes = append(nodes, newNode)
			}
		}

		// Exploración por tags (relacion semántica)
		for _, tag := range currentNote.Tags {
			relatedIDs := g.graph.GetNotesByTag(tag)

			for _, neighborID := range relatedIDs {
				if neighborID == current.id {
					continue
				}

				edgeKey := normalizeEdgeKey(current.id, neighborID)
				if _, exists := edgesSet[edgeKey]; !exists {
					edgesSet[edgeKey] = graphEdge{source: current.id, target: neighborID}
				}

				if !visited[neighborID] && len(nodes) < maxNodes {
					visited[neighborID] = true
					newNode := bfsNode{id: neighborID, depth: current.depth + 1}
					queue = append(queue, newNode)
					nodes = append(nodes, newNode)
				}
			}
		}
	}

	edges := make([]graphEdge, 0, len(edgesSet))
	for _, e := range edgesSet {
		edges = append(edges, e)
	}

	return &bfsResult{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// GetSubgraph construye el subgrafo a partir de una nota.
func (g *GraphService) GetSubgraph(noteID string, depth int, limit int) (*GraphResponse, error) {
	result, err := g.bfs(noteID, depth, limit)
	if err != nil {
		return nil, err
	}

	nodes := make([]GraphNodeDTO, 0, len(result.Nodes))
	nodeSet := make(map[string]bool)

	for _, n := range result.Nodes {
		note, err := g.graph.GetNote(n.id)
		if err != nil {
			continue
		}
		nodes = append(nodes, GraphNodeDTO{
			ID:        note.ID,
			Title:     note.Title,
			IsCurrent: (n.depth == 0),
			Depth:     n.depth,
		})
		nodeSet[note.ID] = true
	}

	// Solo incluir aristas que conectan nodos existentes
	edges := make([]GraphLinkDTO, 0)
	for _, e := range result.Edges {
		if nodeSet[e.source] && nodeSet[e.target] {
			edges = append(edges, GraphLinkDTO{
				Source: e.source,
				Target: e.target,
			})
		}
	}

	return &GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
