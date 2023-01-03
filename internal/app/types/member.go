package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// member hold information of a member
type Member struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty" validate:"required"`
	Name  string             `json:"name" bson:"name" validate:"required"`
	Email string             `json:"email" bson:"email" validate:"omitempty,email"`
}
