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

func GetDistrictsByState(state string) []string {
	values, err := handlers.GetDistinct[string](GetCollection(), "District", bson.M{
		"StateName": state,
	})
	if err != nil {
		return make([]string, 0)
	}
	return values
}

func GetDistrictsByRegion(region string) []string {
	values, err := handlers.GetDistinct[string](GetCollection(), "District", bson.M{
		"RegionName": region,
	})
	if err != nil {
		return make([]string, 0)
	}
	return values
}

func GetRegions(stateName string) []string {
	values, err := handlers.GetDistinct[string](GetCollection(), "RegionName", bson.M{
		"StateName": stateName,
	})
	if err != nil {
		return make([]string, 0)
	}
	return values
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

func getByField(key string, values []string) []DSP {
	docFilter := handlers.DocumentFilter{
		Ctx: context.TODO(),
		Filter: bson.M{
			key: bson.M{
				"$in": values,
			},
		},
	}

	results, err := handlers.GetDocuments[DSP](GetCollection(), docFilter)
	if err != nil {
		return make([]DSP, 0)
	}
	return results
}

func GetByDistrict(districts []string) []DSP {
	return getByField("District", districts)
}

func GetByRegions(regions []string) []DSP {
	return getByField("RegionName", regions)
}

func GetByStates(states []string) []DSP {
	return getByField("StateName", states)
}

func GetByPincodes(pincodes []string) []DSP {
	return getByField("Pincode", pincodes)
}
