package ledgers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sg-business-service/handlers"
)

func GetLedgerByNames(companyId string, ledgerNames []string) []bson.M {
	collectionFilter := bson.M{
		"GUID": bson.M{
			"$regex": "^" + companyId,
		},
		"Name": bson.M{
			"$in": ledgerNames,
		},
	}

	dbFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		Filter:        collectionFilter,
		UsePagination: false,
		Projection: bson.M{
			"Name":        1,
			"Group":       1,
			"Address":     1,
			"MailingName": 1,
			"Email":       1,
			"Emailcc":     1,
		},
	}

	return getCollection().FindDocuments(dbFilter).Data
}
func GetByNames(companyId string, ledgerNames []string) []MetaLedger {
	collectionFilter := bson.M{
		"GUID": bson.M{
			"$regex": "^" + companyId,
		},
		"Name": bson.M{
			"$in": ledgerNames,
		},
	}

	dbFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		Filter:        collectionFilter,
		UsePagination: false,
		Projection: bson.M{
			"Name":        1,
			"Group":       1,
			"Address":     1,
			"MailingName": 1,
			"Email":       1,
			"Emailcc":     1,
		},
	}

	values, err := handlers.GetDocuments[MetaLedger](getCollection(), dbFilter)
	if err != nil {
		return make([]MetaLedger, 0)
	}
	return values
}
