package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectInfo struct {
	UID           primitive.ObjectID `bson:"_id,omitempty"`
	ContentLength int64              `bson:"content_length"`
	ContentType   string             `bson:"content_type"`
	ETag          string             `bson:"etag"`
}
