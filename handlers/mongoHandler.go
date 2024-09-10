package handlers

import (
	"context"
	"fmt"
	"sg-business-service/config"
	"sg-business-service/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentFilter struct {
    Ctx context.Context;
    Filter bson.M;
    UsePagination bool;
    Limit int64;
    Offset int64;
    Projection bson.M;
    Sorting bson.D;
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

func createCollectionIfNotExists(dbName string, collName string) error {
    // Select the database and collection
    collection := Client.Database(dbName).Collection(collName)

    // Perform a write operation to ensure the collection is created
    _, err := collection.InsertOne(context.TODO(), map[string]interface{}{
        "initialized": time.Now(),
    })
    if err != nil {
        return err
    }

    fmt.Printf("Collection '%s' in database '%s' created or already exists.\n", collName, dbName)
    return nil
}


func GetCollection(db string, collection string) *mongo.Collection{
    //createCollectionIfNotExists(db, collection)
    return Client.Database(db).Collection(collection)
}

type MongoHandler struct {
    collection *mongo.Collection
}

func NewMongoHandler(collection *mongo.Collection) *MongoHandler {
    return &MongoHandler{collection: collection}
}


func (handler *MongoHandler) FindDocuments(docFilter DocumentFilter) DocumentResponse {
    findOptions := options.Find()
    if docFilter.UsePagination {
        findOptions.SetSkip(docFilter.Offset * docFilter.Limit)
        findOptions.SetLimit(docFilter.Limit)
    }

    if docFilter.Projection != nil {
        findOptions.SetProjection(docFilter.Projection)
    }
    if docFilter.Sorting != nil {
        findOptions.SetSort(docFilter.Sorting)
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


func GetDocuments[T any](handler *MongoHandler, docFilter DocumentFilter) ([]T, error) {
    findOptions := options.Find()
    if docFilter.UsePagination {
        findOptions.SetSkip(docFilter.Offset * docFilter.Limit)
        findOptions.SetLimit(docFilter.Limit)
    }

    if docFilter.Projection != nil {
        findOptions.SetProjection(docFilter.Projection)
    }
    if docFilter.Sorting != nil {
        findOptions.SetSort(docFilter.Sorting)
    }
    cursor, err := handler.collection.Find(docFilter.Ctx, docFilter.Filter, findOptions)
    if err != nil {
        return nil, err
    }

    defer cursor.Close(docFilter.Ctx)

    var results []T

    for cursor.Next(docFilter.Ctx) {
        var elem T
        if err := cursor.Decode(&elem); err == nil {
            results = append(results, elem)
        }
    }

    if err := cursor.Err(); err != nil {
        return results, err
    }
    return results, err
}


func InsertDocument[T any](dbName string, collName string, document T) (primitive.ObjectID, error) {
    // Select the database and collection
    collection := Client.Database(dbName).Collection(collName)

    // Insert the document
    result, err := collection.InsertOne(context.TODO(), document)
    if err != nil {
        return primitive.NilObjectID, err
    }

    // Get the inserted ID
    insertedID := result.InsertedID.(primitive.ObjectID)
    fmt.Printf("Inserted document with ID: %s\n", insertedID.Hex())
    return insertedID, nil
}

func (handler *MongoHandler) UpdateDocument(dbName string, collectionName string, filter bson.M, update bson.M) {
    // Select the database and collection
    collection := Client.Database(dbName).Collection(collectionName)

    // Update the document
    _, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        utils.ErrorLogger.Fatal(err)
    }
}

func (handler* MongoHandler) AggregatePipeline(dbName string, collectionName string, pipeline mongo.Pipeline) []bson.M {
    // Select the database and collection
    collection := Client.Database(dbName).Collection(collectionName)
    ctx := context.TODO()

    // Execute the aggregation
    cursor, err := collection.Aggregate(ctx, pipeline)

    if err != nil {
        fmt.Println(err)
    }
    defer cursor.Close(ctx)

    var results []bson.M

    for cursor.Next(ctx) {
        var elem bson.M
        if err := cursor.Decode(&elem); err == nil {
            results = append(results, elem)
        }
    }

    if err := cursor.Err(); err != nil {
        fmt.Println(err)
    }
    return results
}

