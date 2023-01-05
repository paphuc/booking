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

// Insert Member to DB Mongo
func (r *MongoRepository) Insert(ctx context.Context, member types.Member) error {
	_, err := r.collection().InsertOne(context.TODO(), member)
	return err
}

// Update Member by using ID
func (r *MongoRepository) UpdateMemberByID(ctx context.Context, member types.Member) error {
	updatedMember := bson.M{"$set": bson.M{
		"name":			member.Name,
		"password":     member.Password,
		"email":		member.Email,
	}}
	_, err := r.collection().UpdateByID(ctx, member.ID, updatedMember)
	return err
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("booking").Collection("members")
}
