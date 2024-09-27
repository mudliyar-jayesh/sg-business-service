package ledgers

import (
	"sg-business-service/models"
)

type MetaLedger struct {
	Name         string                  `bson:"Name"`
	Group        string                  `bson:"Group"`
	Address      *string                 `bson:"Address"`
	State        *string                 `bson:"State"`
	PinCode      *string                 `bson:"PinCode"`
	Email        *string                 `bson:"Email"`
	EmailCc      *string                 `bson:"EmailCc"`
	CreditPeriod *string                 `bson:"CreditPeriod"`
	CreditLimit  *models.FloatFromString `bson:"CreditPeriod"`
}
