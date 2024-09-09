package inventory

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    Name string `bson:"Name"`
}

type ItemGroup struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    Name string `bson:"Name"`
}
