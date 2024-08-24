package endpoints
import (
    "fmt"
    "encoding/json"
    "context"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "sg-business-service/handlers"
)

func GetOutstandingReport(res http.ResponseWriter, req *http.Request) {
    companyId := req.Header.Get("CompanyId")
    var collection = handlers.GetCollection("NewTallyDesktopSync", "Bills")
    var mongoHandler = handlers.NewMongoHandler(collection)

    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter: bson.M{
            "CompanyId": companyId,
        },
        UsePagination: false,
        Limit: int64(0),
        Offset: int64(0),
    }
    var results handlers.DocumentResponse= mongoHandler.FindDocuments(docFilter)

    groupedDocs := GroupByLedgerName(results.Data)

    responseData, err := json.Marshal(groupedDocs)
    if err != nil {
        http.Error(res, "Error encoding response data", http.StatusInternalServerError)
        return
    }

    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(responseData)

}

func GroupByLedgerName(documents []bson.M) map[string][]bson.M {
    grouped := make(map[string][]bson.M)

    for _, doc := range documents {
        ledgerName, exists := doc["LedgerName"].(string)
        if !exists {
            fmt.Println("LedgerName field missing or not a string")
            continue
        }
        grouped[ledgerName] = append(grouped[ledgerName], doc)
    }

    return grouped
}
