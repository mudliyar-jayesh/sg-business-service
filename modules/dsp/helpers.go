package dsp

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func GetStates() []string {
	values, err := handlers.GetDistinct[string](GetCollection(), "StateName", nil)
	if err != nil {
		return make([]string, 0)
	}
	return values
}

func GetPincodes() []string {
	values, err := handlers.GetDistinct[string](GetCollection(), "Pincode", nil)
	if err != nil {
		return make([]string, 0)
	}
	return values
}

func GetDistricts() []string {
	values, err := handlers.GetDistinct[string](GetCollection(), "District", nil)
	if err != nil {
		return make([]string, 0)
	}
	return values
}

func GetByStates(states []string) []DSP {
	docFilter := handlers.DocumentFilter{
		Ctx: context.TODO(),
		Filter: bson.M{
			"StateName": bson.M{
				"$in": states,
			},
		},
	}

	values, err := handlers.GetDocuments[DSP](GetCollection(), docFilter)
	if err != nil {
		return make([]DSP, 0)
	}
	return values
}

func GetByPincode(pincode string) []DSP {
	docFilter := handlers.DocumentFilter{
		Ctx: context.TODO(),
		Filter: bson.M{
			"Pincode": pincode,
		},
	}

	values, err := handlers.GetDocuments[DSP](GetCollection(), docFilter)
	if err != nil {
		return make([]DSP, 0)
	}
	return values
}
