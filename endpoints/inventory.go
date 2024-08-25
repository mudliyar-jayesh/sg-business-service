package endpoints 
import (
    "io"
    "strings"
    "encoding/json"
    "context"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "sg-business-service/handlers"
)

type ItemReportFilter struct {
    StockGroups []string
    SearchText string
    Limit int64
    Offset int64
}


func GetStockItemReport(res http.ResponseWriter, req *http.Request) {
    var collection = handlers.GetCollection("NewTallyDesktopSync", "StockItems")
    var mongoHandler = handlers.NewMongoHandler(collection)

    companyId := req.Header.Get("CompanyId")

    body, err := io.ReadAll(req.Body) 
    if err != nil {
        http.Error(res, "Unable to read request body", http.StatusBadRequest)
        return
    }
    defer req.Body.Close()

    var reqBody ItemReportFilter
    err = json.Unmarshal(body, &reqBody)
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
        tokens := strings.Fields(reqBody.SearchText)

        var regexFilters []bson.M 

        for _, token := range tokens {
            regexFilters  = append(regexFilters, bson.M {
                "Name": bson.M {
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
        UsePagination: reqBody.Limit != 0,
        Limit: reqBody.Limit,
        Offset: reqBody.Offset,
    }

    var results handlers.DocumentResponse = mongoHandler.FindDocuments(docFilter)

    responseData, err := json.Marshal(results.Data)
    if err != nil {
        http.Error(res, "Error encoding response data", http.StatusInternalServerError)
        return
    }


    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(responseData)
}

