package member

import (
	"context"

	"booking/internal/app/types"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository is MongoDB implementation of repository
type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(c *mongo.Client) *MongoRepository {
	return &MongoRepository{
		client: c,
	}
}

// FindByID return member base on given id
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.Member, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var member *types.Member
	err = r.collection().FindOne(ctx, bson.M{"_id": objectID, "disable": false}).Decode(&member)

	return member, err
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("booking").Collection("users")
}
