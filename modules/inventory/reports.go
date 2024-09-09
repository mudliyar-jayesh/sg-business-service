package inventory

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "sg-business-service/handlers"
    "sg-business-service/models"
)

func GetItems(companyId string, pagination models.Pagination, mongoFilter bson.M) []Item {

    var filter = bson.M {
        "GUID": bson.M {
            "$regex": "^"+companyId,
        },
    }
    if mongoFilter != nil {
        filter["$and"] = mongoFilter
    }
    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter:mongoFilter,
        UsePagination: pagination.Apply,
        Limit: pagination.Limit,
        Offset: pagination.Offset,
        Projection: bson.M{
            "Name": 1,
            "_id": 0,
        },
    }

    result, err := handlers.GetDocuments[Item](getItemsCollection(), docFilter)
    if err != nil {
        return make([]Item, 0)
    }
    return result
}

func GetItemGroups(companyId string, pagination models.Pagination, mongoFilter bson.M) []ItemGroup {

    var filter = bson.M {
        "GUID": bson.M {
            "$regex": "^"+companyId,
        },
    }
    if mongoFilter != nil {
        filter["$and"] = mongoFilter
    }
    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter:filter,
        UsePagination: pagination.Apply,
        Limit: pagination.Limit,
        Offset: pagination.Offset,
        Projection: bson.M{
            "Name": 1,
            "_id": 0,
        },
    }

    result, err := handlers.GetDocuments[ItemGroup](getItemGroupsCollection(), docFilter)
    if err != nil {
        return make([]ItemGroup, 0)
    }
    return result
}
