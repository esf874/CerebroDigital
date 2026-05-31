package mongodb

import (
	"context"
	"errors"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoNoteRepository es una implementación de InterfaceNoteRepository usando MongoDB.
type MongoNoteRepository struct {
	collection *mongo.Collection
}

func NewMongoNoteRepository(db *mongo.Database) *MongoNoteRepository {
	return &MongoNoteRepository{
		collection: db.Collection("notes"),
	}
}

// Save - si el id existe, actualiza, sino inserta nuevo documento
func (r *MongoNoteRepository) Save(ctx context.Context, note *domain.Note) error {

	filter := bson.M{"_id": note.ID}
	update := bson.M{"$set": note}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoNoteRepository) FindById(ctx context.Context, id string) (*domain.Note, error) {

	var note domain.Note

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoteNotFound
		}
		return nil, err
	}

	return &note, nil
}

func (r *MongoNoteRepository) Delete(ctx context.Context, id string) error {

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return domain.ErrNoteNotFound
	}

	return nil
}

func (r *MongoNoteRepository) FindAll(ctx context.Context) ([]*domain.Note, error) {
	// bson.M{} es filtro vacío, devuelve todos los documentos
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // cerrar cursor

	var notes []*domain.Note
	for cursor.Next(ctx) {
		var note domain.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

func (r *MongoNoteRepository) FindByTag(ctx context.Context, tag string) ([]*domain.Note, error) {
	filter := bson.M{"tags": tag}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []*domain.Note
	for cursor.Next(ctx) {
		var note domain.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

// SearchByContent busca notas que contengan el query en su contenido o título usando índices de texto de MongoDB.
func (r *MongoNoteRepository) SearchByContent(ctx context.Context, query string) ([]*domain.Note, error) {
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []*domain.Note
	if err = cursor.All(ctx, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}
