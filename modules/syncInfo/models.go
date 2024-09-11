package syncInfo

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type SyncInfo struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    SyncDateTime primitive.DateTime `bson:"SyncDateTime"`
    CompanyId string `bson:"CompanyId"`
}
