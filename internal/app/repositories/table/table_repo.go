package table

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
	return r.client.Database("booking").Collection("tables")
}

// FindByID return member base on given id
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.Table, error) {
	// convert id string to ObjectId
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	var table *types.Table
	err = r.collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&table)
	return table, err
}

// Insert Member to DB Mongo
func (r *MongoRepository) Insert(ctx context.Context, table types.Table) error {
	_, err := r.collection().InsertOne(context.TODO(), table)
	return err
}

// Update Table by using ID
func (r *MongoRepository) UpdateTableByID(ctx context.Context, tableReq types.UpdateTableRequest) error {
	tableId, err := primitive.ObjectIDFromHex(tableReq.ID)
	if err != nil {
		return err
	}

	updatedTable := bson.M{"$set": bson.M{
		"status": tableReq.Status,
	}}

	_, err = r.collection().UpdateByID(ctx, tableId, updatedTable)
	return err
}

// Delete Table by using ID
func (r *MongoRepository) DeleteTable(ctx context.Context, tableReq types.DeleteTableRequest) error {
	tableId, err := primitive.ObjectIDFromHex(tableReq.ID)
	if err != nil {
		return err
	}

	updatedTable := bson.M{"$set": bson.M{
		"del_flg": tableReq.DelFlg,
	}}

	_, err = r.collection().UpdateByID(ctx, tableId, updatedTable)
	return err
}

