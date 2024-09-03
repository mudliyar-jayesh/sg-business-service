package settings

import (
    "sg-business-service/handlers"
)

func GetSettingsCollection() *handlers.MongoHandler {
    settingCollection := handlers.GetCollection("BMRM", "OutstandingSettings")
    return handlers.NewMongoHandler(settingCollection)
}
