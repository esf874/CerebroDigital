package application

import (
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceContainer es una estructura que agrupa todos los servicios de la aplicación, facilitando su inicialización.
type ServiceContainer struct {
	Graph *GraphService
	Note  *NoteService
	Link  *LinkService
	RAG   *RAGService
}

// Crea el grafo compartido, inyecta en servicios, y devuelve el contenedor listo para usar.
func NewServiceContainer(
	noteRepo domain.InterfaceNoteRepository,
	linkRepo domain.InterfaceLinkRepository,
	llmClient LLMClient,
) *ServiceContainer {

	// Crear el grafo único
	sharedGraph := domain.NewGraph()
	graphService := NewGraphService(sharedGraph, noteRepo, linkRepo)

	// Función para inyectar en servicios, modificable en testing
	idGenerator := func() string { return primitive.NewObjectID().Hex() }

	ragService := NewRAGService(sharedGraph, graphService, noteRepo, llmClient)
	linkService := NewLinkService(sharedGraph, linkRepo, noteRepo, idGenerator)

	return &ServiceContainer{
		Graph: graphService,
		Note:  NewNoteService(sharedGraph, noteRepo, linkRepo, linkService, idGenerator),
		Link:  linkService,
		RAG:   ragService,
	}
}
