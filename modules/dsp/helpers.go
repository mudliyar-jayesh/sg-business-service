package dsp

import (
	"sg-business-service/config"
	"sg-business-service/handlers"
)

func GetCollection() *handlers.MongoHandler {
	settingCollection := handlers.GetCollection(config.AppDb, config.DSP)
	return handlers.NewMongoHandler(settingCollection)
}

func InsertOne(entry DSP) {
	handlers.InsertDocument(config.AppDb, config.DSP, entry)
}
