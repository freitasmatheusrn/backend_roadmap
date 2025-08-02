package repo

import (
	"context"

	"github.com/freitasmatheusrn/url-shortening-app/internal/collections"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UrlRepository struct {
	Collection mongo.Collection
}

func NewUrlRepository(db *mongo.Database) *UrlRepository {
	return &UrlRepository{
		Collection: *db.Collection("urls"),
	}
}

func (r *UrlRepository) Create(ctx context.Context, url collections.Url) error {
	_, err := r.Collection.InsertOne(ctx, url)
	return err
}

func (r *UrlRepository) ListAll(ctx context.Context) ([]collections.Url, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []collections.Url
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *UrlRepository) IncrementAccess(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"access_count": 1}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UrlRepository) Update(ctx context.Context, url collections.Url) error {
	filter := bson.M{"_id": url.ID}
	update := bson.M{"$set": bson.M{"original_url": url.OriginalURL}}

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UrlRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	_, err := r.Collection.DeleteOne(ctx, filter)
	return err
}

func (r *UrlRepository) GetOne(ctx context.Context, id string) (string, error){
	var url collections.Url
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&url)
    if err != nil{
		return "Record Not Found", err
	}
	err = r.IncrementAccess(ctx, id)
	if err != nil{
		return "Error updating access count", err
	}
	return url.OriginalURL, nil

}