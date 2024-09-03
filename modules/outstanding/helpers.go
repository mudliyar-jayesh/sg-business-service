package outstanding
import (
    "context"
    "sg-business-service/handlers"
    "go.mongodb.org/mongo-driver/bson"
    "sg-business-service/utils"
)

func getFieldBySortKey(sortKey string) string {
    var fieldBySortKey = make(map[string]string)
    fieldBySortKey["Party"] = "LedgerName"
    fieldBySortKey["Group"] = "LedgerGroupName"
    fieldBySortKey["Bill"] = "Name"

    sortField, exists := fieldBySortKey[sortKey]
    if exists {
        return sortField
    }
    return fieldBySortKey["Party"]
}
func GetFieldBySearchKey(searchKey string) string {
    var fieldBySearchKey = make(map[string]string)
    fieldBySearchKey["Party"] = "LedgerName"
    fieldBySearchKey["Group"] = "LedgerGroupName"
    fieldBySearchKey["Bill"] = "Name"

    searchField, exists := fieldBySearchKey[searchKey]
    if exists {
        return searchField
    }
    return fieldBySearchKey["Party"]
}

func getOutstandingCollection() *handlers.MongoHandler {
    var collection = handlers.GetCollection("NewTallyDesktopSync", "Bills")
    return handlers.NewMongoHandler(collection)
}

func GetOutstandingByFilter(mongoFilter bson.M, requestFilter OsReportFilter, usePagination bool) handlers.DocumentResponse {
    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter:mongoFilter,
        UsePagination: usePagination,
        Limit: requestFilter.Limit,
        Offset: requestFilter.Offset,
        Projection: bson.M{
            "LedgerName": 1,
            "LedgerGroupName": 1,
            "BillDate": "$BillDate.Date",
            "DueDate": "$BillCreditPeriod.DueDate",
            "Amount": "$ClosingBal.Amount",
            "Name": "$Name",
            "_id": 0,
        },
        Sorting: bson.D {
            {
                Key: getFieldBySortKey(requestFilter.SortKey),
                Value: utils.GetValueBySortOrder(requestFilter.SortOrder),
            },
        },
    }

    var handler  = getOutstandingCollection()
    return handler.FindDocuments(docFilter)
}

