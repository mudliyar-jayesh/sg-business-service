package syncInfo

import (
    "context"
    "sg-business-service/handlers"
    "go.mongodb.org/mongo-driver/bson"
)

func GetByCompanyId(companyId string) *SyncInfo {
    collection := GetCollection()

    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        UsePagination: false,
        Filter:  bson.M {
            "CompanyId" : companyId,
        },
        Projection: bson.M {
            "_id": 0,
            "CompanyId": 1,
            "SyncDateTime": 1,
        },
    }

    result, err := handlers.GetDocuments[SyncInfo](collection, docFilter)
    if err != nil || len(result) < 1 {
        return nil
    }

    return &result[0]
}
