package ledgers 

import (
    "sg-business-service/handlers"
    "sg-business-service/config"
)

func getCollection() *handlers.MongoHandler {
    var collection = handlers.GetCollection(config.TallyDb, config.Ledger)
    return handlers.NewMongoHandler(collection)
}

