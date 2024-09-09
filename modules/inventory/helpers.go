package inventory

import (
    "sg-business-service/handlers"
    "sg-business-service/config"
)

func getItemsCollection() *handlers.MongoHandler {
    var collection = handlers.GetCollection(config.TallyDb, config.Item)
    return handlers.NewMongoHandler(collection)
}

func getItemGroupsCollection() *handlers.MongoHandler {
    var collection = handlers.GetCollection(config.TallyDb, config.ItemGroup)
    return handlers.NewMongoHandler(collection)
}


