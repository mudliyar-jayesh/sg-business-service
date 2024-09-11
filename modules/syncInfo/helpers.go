package syncInfo

import (
    "sg-business-service/config"
    "sg-business-service/handlers"
)

func GetCollection() *handlers.MongoHandler {
    settingCollection := handlers.GetCollection(config.TallyDb, config.SyncInfo)
    return handlers.NewMongoHandler(settingCollection)
}

