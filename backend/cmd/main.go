package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/api/handlers"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/application"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/infrastructure/llm"
	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/infrastructure/persistence/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
)

const (
	urlNote = "/api/notes/{id}"
)

func main() {

	// Carga variables de entorno del .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env")
	}

	// lee variables de entorno
	uri := os.Getenv("MONGO_URI")
	nameDB := os.Getenv("NAME_DATABASE")
	port := os.Getenv("APP_PORT")
	llmURL := os.Getenv("LLM_URL")

	if uri == "" || nameDB == "" {
		log.Fatal("Faltan variables de entorno")
	}
	if llmURL == "" {
		llmURL = "http://localhost:8081"
	}

	// backend se conecta a llama-server puerto 8081
	llmClient := llm.NewClient(llmURL)

	// crear conexion mongo
	client, err := mongodb.NewMongoClient(uri)
	if err != nil {
		log.Fatal("Error conectando a Mongo:", err)
	}
	fmt.Println("Conexión establecida a MongoDB.")

	// poner timeout a la desconexion cuando acabe
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = client.Disconnect(ctx); err != nil {
			log.Println("Error al desconectar MongoDB:", err)
		} else {
			fmt.Println("Desconexión realizada correctamente.")
		}
	}()

	// obtener bd + crear indices
	db := client.Database(nameDB)
	ctx := context.Background()
	noteCollection := db.Collection("notes")
	linkCollection := db.Collection("links")

	// indice de texto para searchByContent
	if _, err := noteCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "content", Value: "text"},
			{Key: "title", Value: "text"},
		},
	}); err != nil {
		log.Fatal("Error crenado índice de búsqueda en notes:", err)
	}

	// indices para consultas frecuentes de links por id
	_, err = linkCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "origin_note_id", Value: 1}}},
		{Keys: bson.D{{Key: "dest_note_id", Value: 1}}},
	})
	if err != nil {
		log.Fatal("Error creando índices en links:", err)
	}

	// Crear repos + inyectar bd
	noteRepo := mongodb.NewMongoNoteRepository(db)
	linkRepo := mongodb.NewMongoLinkRepository(db)

	// Crear contenedor de servicios, el constructor recibe interfaces de los repos, se pasa implem
	// el container crea servicios (note, link, graph) integrando grafo compartido y los repos
	container := application.NewServiceContainer(noteRepo, linkRepo, llmClient)

	// cargar el grafo desde BD al arrancar la app
	// reconstruye en memoria el grafo a partir de las notas y links de la bd
	if err := container.Graph.LoadGraph(ctx); err != nil {
		log.Println("Error cargando el grafo desde la BD:", err)
	}

	// Crear handlers y ruta
	ragHandler := handlers.NewRAGHandler(container.RAG)
	noteHandler := handlers.NewNoteHandler(container.Note)
	graphHandler := handlers.NewGraphHandler(container.Graph)

	// rutas
	router := mux.NewRouter()

	router.HandleFunc("/api/ask", ragHandler.Ask).Methods("POST")
	router.HandleFunc("/api/notes", noteHandler.GetAllNotes).Methods("GET")
	router.HandleFunc("/api/notes", noteHandler.CreateNote).Methods("POST")

	// Rutas con parámetros
	router.HandleFunc("/api/notes/lookup", noteHandler.GetNoteByTitle).Queries("title", "{title}").Methods("GET")
	router.HandleFunc(urlNote, noteHandler.GetNote).Methods("GET")
	router.HandleFunc(urlNote, noteHandler.UpdateNote).Methods("PUT")
	router.HandleFunc(urlNote, noteHandler.DeleteNote).Methods("DELETE")

	router.HandleFunc("/api/notes/{id}/graph", graphHandler.GetGraph).Methods("GET")

	router.HandleFunc("/api/notes/{id}/tags/{tag}", noteHandler.AddTag).Methods("POST")
	router.HandleFunc("/api/notes/{id}/tags/{tag}", noteHandler.RemoveTag).Methods("DELETE")

	// iniciar servidor
	log.Printf("Servidor escuchando en :%s", port)
	http.ListenAndServe(":"+port, router)
}
