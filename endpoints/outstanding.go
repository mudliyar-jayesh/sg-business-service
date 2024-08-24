package endpoints
import (
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

    res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(res).Encode(results.Data)
}
