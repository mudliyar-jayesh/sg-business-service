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
    "strconv"
)

type OsReportFilter struct {
    PartyName string
    SearchText string
    Limit int64
    Offset int64
    Groups []string
    OnlyPending bool
    OnlyDue bool
    OnlyOverDue bool
    DueDays int
    OverDueDays int
}

func GetCachedGroups(res http.ResponseWriter, req *http.Request) {
    companyId := req.Header.Get("CompanyId")

    isDebitStr := req.URL.Query().Get("isDebit")
    if isDebitStr == "" {
        isDebitStr = "true"
    }
    var isDebit bool
    isDebit, _ = strconv.ParseBool(isDebitStr)

    var parentName string = "Current Assets"
    if !isDebit {
        parentName = "Current Liabilities"
    }

    var groups = handlers.CachedGroups.GetChildrenNames(companyId, parentName)
    responseData, err := json.Marshal(groups)
    if err != nil {
        http.Error(res, "Error encoding response data", http.StatusInternalServerError)
        return
    }

    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(responseData)

}


func GetOutstandingReport(res http.ResponseWriter, req *http.Request) {
    var collection = handlers.GetCollection("NewTallyDesktopSync", "Bills")
    var mongoHandler = handlers.NewMongoHandler(collection)

    companyId := req.Header.Get("CompanyId")
    isDebitStr := req.URL.Query().Get("isDebit")
    if isDebitStr == "" {
        isDebitStr = "true"
    }
    var isDebit bool
    isDebit, _ = strconv.ParseBool(isDebitStr)

    var parentName string = "Current Assets"
    if !isDebit {
        parentName = "Current Liabilities"
    }

    var groups = handlers.CachedGroups.GetChildrenNames(companyId, parentName)

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

    if len(reqBody.Groups) > 0 {
        groups = utils.Intersection(groups, reqBody.Groups)
    }

    var filter = bson.M {
        "CompanyId": companyId,
        "LedgerGroupName": bson.M {
            "$in": groups,
        },
        "ClosingBal": bson.M {"$ne": nil },
        "ClosingBal.IsDebit": isDebit,
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
        Projection: bson.M{
            "LedgerName": 1,
            "LedgerGroupName": 1,
            "BillDate.Date": 1,
            "BillCreditPeriod.DueDate": 1,
            "ClosingBal.Amount": 1,
            "_id": 0,
        },
    }

    var results handlers.DocumentResponse= mongoHandler.FindDocuments(docFilter)

    /*
    var bills []interface{}
    for _, data := range results.Data {

        var dueDateStr string

        if billCreditPeriod, ok := data["BillCreditPeriod"].(map[string]interface{}); ok {
            if dueDate, ok := billCreditPeriod["DueDate"].(string); ok && dueDate != "" {
                dueDateStr = dueDate
            }
        }

        if dueDateStr == "" {
            if billDate, ok := data["BillDate"].(map[string]interface{}); ok {
                if date, ok := billDate["Date"].(string); ok && date != "" {
                    dueDateStr = date
                }
            }
        }

        if dueDateStr == "" {
            fmt.Println("No due date found")
            return
        }

        // Parse the DueDate
        dueDate, err := time.Parse(time.RFC3339, dueDateStr)
        if err != nil {
            fmt.Println("Error parsing DueDate:", err)
            return
        }

        // Get today's date
        today := time.Now().UTC()

        // Calculate the difference in days
        daysUntilDue := int(dueDate.Sub(today).Hours() / 24)

        var overview = OsOverview{
            PartyName: data["LedgerName"].(string),
        }
        if daysUntilDue < reqBody.DueDays {
        }
 
        bills = append(bills, overview)
    }

    */
    var temp = Temp {
        Data:  results.Data,
        Count: len(results.Data),
    }

    responseData, err := json.Marshal(temp)
    if err != nil {
        http.Error(res, "Error encoding response data", http.StatusInternalServerError)
        return
    }


    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(responseData)

}

type Temp struct {
    Data interface{}
    Count int
}
type OsOverview struct {
    PartyName string
    PendingAmount float64
    OverDueAmount float64
    DueAmount float64
}
