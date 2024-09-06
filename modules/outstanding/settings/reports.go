package settings
import (
    "context"
    "sg-business-service/handlers"
    "go.mongodb.org/mongo-driver/bson"
    "sg-business-service/models"
)

func GetAllSettings(companyId string) ([]models.OsShareSettings, error) {
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
    return handlers.GetDocuments[models.OsShareSettings](handler, docFilter)
}


