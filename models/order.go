package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductId primitive.ObjectID `json:"product_id" bson:"product_id"`
	Price     int                `json:"price" bson:"price"`
	Quantity  int                `json:"quantity" bson:"quantity"`
}

type Order struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Cart       []Cart             `json:"cart" bson:"cart"`
	TotalPrice int                `json:"total_price" bson:"total_price"`
	Address    string             `json:"address" bson:"address"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}
