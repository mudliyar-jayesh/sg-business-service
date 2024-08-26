package endpoints
import (
    "time"
    "strings"
    "encoding/json"
    "context"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
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

    response := utils.NewResponseStruct(responseData, len(responseData))
    response.ToJson(res)
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

    reqBody, err := utils.ReadRequestBody[OsReportFilter](req)
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
        UsePagination: reqBody.Limit != 0,
        Limit: reqBody.Limit,
        Offset: reqBody.Offset,
        Projection: bson.M{
            "LedgerName": 1,
            "LedgerGroupName": 1,
            "BillDate": "$BillDate.Date",
            "DueDate": "$BillCreditPeriod.DueDate",
            "Amount": "$ClosingBal.Amount",
            "Name": "$Name",
            "_id": 0,
        },
    }

    var results handlers.DocumentResponse= mongoHandler.FindDocuments(docFilter)

    var bills []Bill
    istLocation, _ := time.LoadLocation("Asia/Kolkata")
    for _, item := range results.Data {
        billDateValue := item["BillDate"].(primitive.DateTime).Time()
        billDate := billDateValue.In(istLocation).Format("2006-01-02 15:04:05")

        var dueDate string
        dueDateValue := item["DueDate"]
        if dueDateValue == nil {
            dueDate = billDate
        } else {
           dueDate = dueDateValue.(primitive.DateTime).Time().In(istLocation).Format("2006-01-02 15:04:05")
        }

        layout := "2006-01-02 15:04:05"
        parsedTime, _:= time.Parse(layout, dueDate)
        today := time.Now().UTC()

        // Calculate the difference
        diff := today.Sub(parsedTime)

        // Get the difference in days
        days := int(diff.Hours() / 24)

        var bill Bill = Bill {
            LedgerName: item["LedgerName"].(string),
            LedgerGroupName: item["LedgerGroupName"].(string),
            BillName: item["Name"].(string),
            DueDate: dueDate,
            BillDate: billDate,
            DelayDays: days,
        }
        var amount = parseFloat64(item["Amount"])
        if days >= reqBody.DueDays && days <= reqBody.OverDueDays {
            bill.DueAmount = amount
        } else if days > reqBody.OverDueDays {
            bill.OverDueAmount = amount
        } else {
            bill.Amount = amount
        }

        bills = append(bills, bill)
    }

    var groupedBills = utils.GroupByKey[Bill](bills, "LedgerName")

    response := utils.NewResponseStruct(groupedBills, len(groupedBills))
    response.ToJson(res)
}
func parseFloat64(value interface{}) float64 {
    var result float64
    switch v := value.(type) {
    case float64:
        result = v
    case int:
        result = float64(v)
    case string:
        parsed, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return 0 // Return default value on error
        }
        result = parsed
    default:
        return 0 // Return default value if type is not handled
    }
    return result
}

type Temp struct {
    Data interface{}
    Count int
}
type Bill struct {
    LedgerName string 
    LedgerGroupName string 
    BillName string
    BillDate string
    DueDate string
    DelayDays int
    Amount float64
    DueAmount float64
    OverDueAmount float64
}
