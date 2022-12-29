package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	FileID        primitive.ObjectID `bson:"_id"`
	User          primitive.ObjectID `bson:"user"`
	FileName      string             `bson:"fileName"`
	FileExtension string             `bson:"fileExtension"`
	FileLocation  string             `bson:"fileLocation"`
	InsertionDate time.Time          `bson:"insertionDate"`
}
