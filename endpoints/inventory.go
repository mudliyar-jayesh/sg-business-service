package endpoints 
import (
    "context"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "sg-business-service/handlers"
    "sg-business-service/utils"
    "sg-business-service/models"
    "sg-business-service/modules/inventory"
)

type ItemReportFilter struct {
    StockGroups []string
    SearchText string
    SearchKey string
    SortKey string;
    SortOrder string;
    Limit int64
    Offset int64
}

func getItemFieldBySortKey(sortKey string) string {
    var fieldBySortKey = make(map[string]string)
    fieldBySortKey["Item"] = "Name"
    fieldBySortKey["Group"] = "StockGroup"
    fieldBySortKey["Quantity"] = "ClosingBal.Number"
    fieldBySortKey["Rate"] = "ClosingRate.RatePerUnit"
    fieldBySortKey["Amount"] = "ClosingValue.Amount"

    sortField, exists := fieldBySortKey[sortKey]
    if exists {
        return sortField
    }
    return fieldBySortKey["Item"]
}
func getItemFieldBySearchKey(searchKey string) string {
    var fieldBySearchKey = make(map[string]string)
    fieldBySearchKey["Item"] = "Name"
    fieldBySearchKey["Group"] = "StockGroup"

    searchField, exists := fieldBySearchKey[searchKey]
    if exists {
        return searchField
    }
    return fieldBySearchKey["Item"]
}


func GetStockItemReport(res http.ResponseWriter, req *http.Request) {
    var collection = handlers.GetCollection("NewTallyDesktopSync", "StockItems")
    var mongoHandler = handlers.NewMongoHandler(collection)

    companyId := req.Header.Get("CompanyId")

    reqBody, err := utils.ReadRequestBody[ItemReportFilter](req)
    if err != nil {
        http.Error(res, "Unable to read request body", http.StatusBadRequest)
        return
    }

    var filter = bson.M {
        "GUID": bson.M {
            "$regex": "^"+companyId,
        },
    }
    if len(reqBody.StockGroups) > 0 {
        filter["StockGroup"] = bson.M {
            "$in": reqBody.StockGroups,
        }
    }
    if len(reqBody.SearchText) > 0 {
        var searchField = getItemFieldBySearchKey(reqBody.SearchKey)
        filter["$and"] = utils.GenerateSearchFilter(reqBody.SearchText, searchField)
    }
    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter:filter,
        UsePagination: reqBody.Limit != 0,
        Limit: reqBody.Limit,
        Offset: reqBody.Offset,
        Sorting: bson.D {
            {
                Key: getItemFieldBySortKey(reqBody.SortKey),
                Value: utils.GetValueBySortOrder(reqBody.SortOrder),
            },
        },
    }

    var results handlers.DocumentResponse = mongoHandler.FindDocuments(docFilter)

    response := utils.NewResponseStruct(results.Data, len(results.Data))
    response.ToJson(res)
}

func GetItemGroupNames(res http.ResponseWriter, req *http.Request) {
    companyId := req.Header.Get("CompanyId")
    searchKey := req.URL.Query().Get("searchKey")

    var filter bson.M = nil
    var pagination = models.Pagination {
        Apply: true,
        Limit: 25,
        Offset: 0,
    }

    if len(searchKey) > 0 {
        pagination.Apply = false
        filter = bson.M {
        }
        filter["$and"] = utils.GenerateSearchFilter(searchKey, "Name")
    }
    itemGroups := inventory.GetItemGroups(companyId, pagination, filter)
    response := utils.NewResponseStruct(itemGroups, len(itemGroups))
    response.ToJson(res)

}
