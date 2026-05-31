package mongodb

import (
	"context"

	"gitlab.com/HP-SCDS/Observatorio/2025-2026/cerebrodigital/usal-za-cerebrodigital/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoLinkRepository es una implementación de InterfaceLinkRepository usando MongoDB.
type MongoLinkRepository struct {
	collection *mongo.Collection
}

func NewMongoLinkRepository(db *mongo.Database) *MongoLinkRepository {
	return &MongoLinkRepository{
		collection: db.Collection("links"),
	}
}

func (r *MongoLinkRepository) Save(ctx context.Context, link *domain.Link) error {
	filter := bson.M{"_id": link.ID}
	update := bson.M{"$set": link}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoLinkRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrLinkNotFound
	}
	return nil
}

func (r *MongoLinkRepository) DeleteAllByNoteID(ctx context.Context, noteID string) error {
	filter := bson.M{
		"$or": []bson.M{
			{"origin_note_id": noteID},
			{"dest_note_id": noteID},
		},
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

func (r *MongoLinkRepository) FindByOrigin(ctx context.Context, noteID string) ([]*domain.Link, error) {
	return r.find(ctx, bson.M{"origin_note_id": noteID})
}

func (r *MongoLinkRepository) FindByDest(ctx context.Context, noteID string) ([]*domain.Link, error) {
	return r.find(ctx, bson.M{"dest_note_id": noteID})
}

func (r *MongoLinkRepository) FindByNoteID(ctx context.Context, noteID string) ([]*domain.Link, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"origin_note_id": noteID},
			{"dest_note_id": noteID},
		},
	}
	return r.find(ctx, filter)
}

func (r *MongoLinkRepository) FindAll(ctx context.Context) ([]*domain.Link, error) {
	return r.find(ctx, bson.M{})
}

// Método auxiliar para no repetir codigo (recibe filtro)
func (r *MongoLinkRepository) find(ctx context.Context, filter bson.M) ([]*domain.Link, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var links []*domain.Link
	if err = cursor.All(ctx, &links); err != nil {
		return nil, err
	}
	return links, nil
}
