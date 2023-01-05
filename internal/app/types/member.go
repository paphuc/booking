package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// member hold information of a member
type Member struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty" validate:"required"`
	Name  string             `json:"name" bson:"name" validate:"required"`
	Password string 		 `json:"password" bson:"password" validate:"required"`
	Email string             `json:"email" bson:"email" validate:"required"`
}

type MemberRequest struct {
	Name  string             `json:"name" validate:"required"`
	Password string 		 `json:"password" validate:"required"`
	Email    string 		 `json:"email" validate:"required,email"`
}

type MemberResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
