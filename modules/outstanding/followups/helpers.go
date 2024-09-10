package followups

import (
	"context"
	"fmt"
	"sg-business-service/config"
	"sg-business-service/handlers"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// ------------ FOLLOW UP ----------------
func getFollowupCollection() *handlers.MongoHandler {
    var collection = handlers.GetCollection(config.AppDb, config.FollowUp)
    return handlers.NewMongoHandler(collection)
}

func getFollowupById (companyId, id string) *FollowUp{
	filter := handlers.DocumentFilter{UsePagination: false, Ctx: context.TODO(), Filter: bson.M{
		"CompanyId": companyId,
		"FollowUpId": id,
	}}

	collection := getFollowupCollection()

	res, err := handlers.GetDocuments[FollowUp](collection, filter)

	if err != nil || len(res) < 1{
		return nil
	}

	return &res[0]
}
// ------------ END FOLLOW UP ----------------


// ------------ CONTACT PERSON COLLECTION ------------
func getContactCollection() *handlers.MongoHandler {
    var collection = handlers.GetCollection(config.AppDb, config.ContactPerson)
    return handlers.NewMongoHandler(collection)
}

func getContactPersonById (companyId, id string) *ContactPerson{
	filter := handlers.DocumentFilter{UsePagination: false, Ctx: context.TODO(), Filter: bson.M{
		"CompanyId": companyId,
		"PersonId": id,
	}}

	collection := getContactCollection() 

	res, err := handlers.GetDocuments[ContactPerson](collection, filter)

	if err != nil || len(res) < 1{
		return nil
	}

	return &res[0]
}

func createContactPerson (person ContactPerson)  {
	newUUid := uuid.New()
	person.PersonId = newUUid.String()

	_, err := handlers.InsertDocument[ContactPerson](config.AppDb, config.ContactPerson, person)

	if err != nil {
		fmt.Println(err);
	}
}


func insertFollowUpToDB(followup FollowUp) (string, error) {
	guid := uuid.New()
	followup.FollowUpId = guid.String() 

	_, err := handlers.InsertDocument[FollowUp](config.AppDb, config.FollowUp, followup)

	if err != nil {
		fmt.Println(err);
	}

	return followup.FollowUpId, err
}
// ------------ END CONTACT PERSON COLLECTION -----------


