package settings

import (
	"context"
	"sg-business-service/handlers"
	"sg-business-service/modules/outstanding"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllSettings(companyId string) ([]outstanding.OsShareSettings, error) {
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
    return handlers.GetDocuments[outstanding.OsShareSettings](handler, docFilter)
}


