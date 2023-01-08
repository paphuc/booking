package types

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"	
)

// Table hold information of a table
type Table struct {
	ID    			primitive.ObjectID		`json:"_id" bson:"_id,omitempty" validate:"required"`
	Status			string					`json:"status" bson:"status" validate:"required"`
	Type			string					`json:"type" bson:"type" validate:"required"`
	Slots 			int 					`json:"slots" bson:"slots" validate:"required,max=10"`
	DelFlg 			bool 					`json:"del_flg" bson:"del_flg" validate:"omitempty"`
	CreateAt     	time.Time         		`json:"create_at" bson:"create_at"`
	UpdateAt     	time.Time          		`json:"update_at" bson:"update_at"`
	// RestaurantId 	primitive.ObjectID   	`json:"restaurant_id" bson:"restaurant_id validate:"required"` 
}

type TableRequest struct {
	Status			string					`json:"status" bson:"status" validate:"required"`
	Type			string					`json:"type" bson:"type" validate:"required"`
	Slots 			int 					`json:"slots" bson:"slots" validate:"required,max=10"`
}

type UpdateTableRequest struct {
	ID       		string 					`json:"_id" bson:"_id,omitempty" validate:"required"`
    Status	 		string					`json:"status" bson:"status" validate:"required"`
}

type DeleteTableRequest struct {
	ID       		string 					`json:"_id" bson:"_id,omitempty" validate:"required"`
    DelFlg 			bool 					`json:"del_flg" bson:"del_flg" validate:"required"`
}
