package settings
import (
    "context"
    "sg-business-service/handlers"
    "go.mongodb.org/mongo-driver/bson"
)

func GetAllSettings(companyId string) handlers.DocumentResponse {
    docFilter  := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter: bson.M {
            "CompanyId": companyId,
        },
        UsePagination: false,
        Limit: 0,
        Offset: 0,
    }

    var handler = GetSettingsCollection()
    return handler.FindDocuments(docFilter)
}


