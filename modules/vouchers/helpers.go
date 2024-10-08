package vouchers

import (
  "sg-business-service/config"
	"sg-business-service/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *handlers.MongoHandler {
	var collection = handlers.GetCollection(config.TallyDb, config.Bill)
	return handlers.NewMongoHandler(collection)
}


// This function will execute the aggregation pipeline as per your provided C# code.
func GetMetaVouchers(companyId string, voucherTypes []string, voucherIds []string, ledgerNames []string) []MetaVoucher {
	pipeline := mongo.Pipeline{}

	// Match part for GUID and VoucherType
	match := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "GUID", Value: bson.D{{Key: "$regex", Value: companyId}}},
			{Key: "VoucherType", Value: bson.D{{Key: "$in", Value: voucherTypes}}},
		}},
	}

	// Check if voucherIds are provided and apply the filter
	if len(voucherIds) > 0 {
		match = bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "GUID", Value: bson.D{{Key: "$in", Value: voucherIds}}},
			}},
		}
	}

	// Append the match stage to the pipeline
	pipeline = append(pipeline, match)

	// Apply LedgerName filter if provided
	if len(ledgerNames) > 0 {
		ledgerNameMatch := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$expr", Value: bson.D{
					{Key: "$in", Value: bson.A{
						bson.D{
							{Key: "$arrayElemAt", Value: bson.A{"$Ledgers.LedgerName", 0}},
						},
						ledgerNames,
					}},
				}},
			}},
		}

		// Append the ledgerName match stage to the pipeline
		pipeline = append(pipeline, ledgerNameMatch)
	}

	// Project stage to select the required fields
project := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "GUID", Value: 1},
			{Key: "Date", Value: "$Date.Date"},
			{Key: "VoucherType", Value: 1},
			{Key: "ReferenceDate", Value: "$ReferenceDate.Date"},
			{Key: "Reference", Value: 1},
			{Key: "InventoryAllocations", Value: 1},
			{Key: "Ledgers", Value: 1},
		}},
	// Append the project stage to the pipeline
  }
	pipeline = append(pipeline, project)

	// Execute the aggregation pipeline
  result, err := handlers.AggregateCollection[MetaVoucher](config.TallyDb, config.Voucher, pipeline)
  if err != nil {
    return make([]MetaVoucher, 0)
  }

	return result
}


