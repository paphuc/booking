package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// restaurant hold information of a restaurant
type Restaurant struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty" validate:"required"`
	Name  string             `json:"name" bson:"name" validate:"required"`
	Address string 		     `json:"address" bson:"address" validate:"omitempty,address"`
}
