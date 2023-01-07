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

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("booking").Collection("members")
}

// FindByID return member base on given id
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.Member, error) {
	// convert id string to ObjectId
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	var member *types.Member
	err = r.collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&member)
	return member, err
}

// Insert Member to DB Mongo
func (r *MongoRepository) Insert(ctx context.Context, member types.Member) error {
	_, err := r.collection().InsertOne(context.TODO(), member)
	return err
}

// Update Member by using ID
func (r *MongoRepository) UpdateMemberByID(ctx context.Context, member types.UpdateMemberRequest) error {
	memberId, err := primitive.ObjectIDFromHex(member.ID)
	if err != nil {
		return err
	}

	updatedMember := bson.M{"$set": bson.M{
		"password": member.Password,
	}}

	_, err = r.collection().UpdateByID(ctx, memberId, updatedMember)
	return err
}

func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (*types.Member, error) {
	var user *types.Member
	err := r.collection().FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}
