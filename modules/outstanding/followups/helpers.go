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

func getFollowupListByParty (companyId, partyName string) []FollowUp{
	filter := handlers.DocumentFilter{UsePagination: false, Ctx: context.TODO(), Filter: bson.M{
		"CompanyId": companyId,
		"PartyName": partyName,
	}}

	collection := getFollowupCollection()

	res, err := handlers.GetDocuments[FollowUp](collection, filter)

	if err != nil || len(res) < 1{
		return nil
	}

	return res
}

func updateFollowup(followup FollowUp) error {
	mongoHandler :=getFollowupCollection()

	filter := bson.M {
		"FollowUpId":followup.FollowUpId, 
	}

	err := mongoHandler.ReplaceDocument(config.AppDb, config.FollowUp, filter, followup)

	return err
}

func insertFollowUpToDB(followup FollowUp) (string, error) {
	guid := uuid.New()
	followup.FollowUpId = guid.String() 

	id, err := handlers.InsertDocument[FollowUp](config.AppDb, config.FollowUp, followup)

	if err != nil {
		fmt.Println(err);
	}else{
		fmt.Println(id.String())	
	}

	return followup.FollowUpId, err
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

func getContactPersonList (companyId string, partyName string) []ContactPerson{
	filter := handlers.DocumentFilter{UsePagination: false, Ctx: context.TODO(), Filter: bson.M{
		"CompanyId": companyId,
		"PartyName": partyName,
	}}

	collection := getContactCollection() 

	res, err := handlers.GetDocuments[ContactPerson](collection, filter)

	if err != nil || len(res) < 1{
		return nil
	}

	return res
}

func createContactPerson (person ContactPerson) (string, error) {
	guid := uuid.New()
	person.PersonId = guid.String()

	_, err := handlers.InsertDocument[ContactPerson](config.AppDb, config.ContactPerson, person)

	return person.PersonId, err
}
// ------------ END CONTACT PERSON COLLECTION -----------


