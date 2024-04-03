package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	UID         primitive.ObjectID `bson:"_id,omitempty"`
	ID          int                `bson:"id" `
	Title       string             `bson:"title"`
	Price       float64            `bson:"price"`
	Category    string             `bson:"category"`
	Brand       string             `bson:"brand" `
	Url         string             `bson:"url" `
	Description string             `bson:"description"`
}
