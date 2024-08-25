package endpoints
import (
    "io"
    "strings"
    "encoding/json"
    "context"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "sg-business-service/handlers"
    "sg-business-service/utils"
)

type OsReportFilter struct {
    PartyName string
    SearchText string
    Limit int64
    Offset int64
}


func GetOutstandingReport(res http.ResponseWriter, req *http.Request) {
    companyId := req.Header.Get("CompanyId")
    var collection = handlers.GetCollection("NewTallyDesktopSync", "Bills")
    var mongoHandler = handlers.NewMongoHandler(collection)


    body, err := io.ReadAll(req.Body) 
    if err != nil {
        http.Error(res, "Unable to read request body", http.StatusBadRequest)
        return
    }
    defer req.Body.Close()

    var reqBody OsReportFilter
    err = json.Unmarshal(body, &reqBody)
    if err != nil {
        http.Error(res, "Unable to read request body", http.StatusBadRequest)
        return
    }

    var filter = bson.M {
        "CompanyId": companyId,
    }

    if len(reqBody.PartyName) > 0 {
        filter["LedgerName"] = reqBody.PartyName
    } else if len(reqBody.SearchText) > 0 {
        tokens := strings.Fields(reqBody.SearchText)

        var regexFilters []bson.M 

        for _, token := range tokens {
            regexFilters  = append(regexFilters, bson.M {
                "LedgerName": bson.M {
                    "$regex": token,
                    "$options": "i", //case insensitive
                },
            })
        }

        filter["$and"] = regexFilters
    }


    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter:filter,
        UsePagination: reqBody.Limit != 0 && reqBody.Offset != 0,
        Limit: reqBody.Limit,
        Offset: reqBody.Offset,
    }
    var results handlers.DocumentResponse= mongoHandler.FindDocuments(docFilter)

    groupedDocs := utils.GroupBy(results.Data, "LedgerName")

    responseData, err := json.Marshal(groupedDocs)
    if err != nil {
        http.Error(res, "Error encoding response data", http.StatusInternalServerError)
        return
    }

    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(responseData)

}
