package handlers

import (
    "fmt"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "sg-business-service/utils"
    "sg-business-service/config"
)

type DocumentFilter struct {
    Ctx context.Context;
    Filter bson.M;
    UsePagination bool;
    Limit int64;
    Offset int64;
}

type DocumentResponse struct {
    Data []bson.M;
    Err error;
}

var Client *mongo.Client

func ConnectToMongo(cfg *config.MongoConfig) {
    clientOptions := options.Client().ApplyURI(cfg.Uri)
    var err error 
    Client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        utils.ErrorLogger.Fatal("Failed to connect to MongoDB: ", err)
    }
    if err := Client.Ping(context.TODO(), nil); err != nil {
        utils.ErrorLogger.Fatal("Failed to ping MongoDB: ", err)
    }
    fmt.Println("[*]Connected to Mongo")
}

func GetCollection(db string, collection string) *mongo.Collection{
    return Client.Database(db).Collection(collection)
}

type MongoHandler struct {
    collection *mongo.Collection
}

func NewMongoHandler(collection *mongo.Collection) *MongoHandler {
    return &MongoHandler{collection: collection}
}


func (handler *MongoHandler) FindDocuments(docFilter DocumentFilter) DocumentResponse {
    var findOptions *options.FindOptions
    if docFilter.UsePagination {
        findOptions = &options.FindOptions{
            Skip: &docFilter.Offset,
            Limit: &docFilter.Limit,
        }
    }
    cursor, err := handler.collection.Find(docFilter.Ctx, docFilter.Filter, findOptions)
    if err != nil {
        return DocumentResponse {
            Data: nil,
            Err: err,
        }
    }

    defer cursor.Close(docFilter.Ctx)

    var results []bson.M

    for cursor.Next(docFilter.Ctx) {
        var elem bson.M
        if err := cursor.Decode(&elem); err == nil {
            results = append(results, elem)
        }
    }

    if err := cursor.Err(); err != nil {
        return DocumentResponse {
            Data: nil,
            Err: err,
        }
    }
    return DocumentResponse {
        Data: results,
        Err: err,
    }
}
