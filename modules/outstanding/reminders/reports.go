package reminders

import (
    "fmt"
    "time"
    "context"
    "sg-business-service/handlers"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "sg-business-service/utils"
    osMod "sg-business-service/modules/outstanding"
    //configMod "sg-business-service/modules/outstanding/settings"
    ledgersMod "sg-business-service/modules/ledgers"
    "sg-business-service/models"
)

func SendEmailReminder(companyId string, ledgerNames []string) {
    ledgers := ledgersMod.GetLedgerByNames(companyId, ledgerNames)
    ledgerByName, err := utils.ToDictionary(ledgers, "Name")
    if err != nil {
        return
    }

    //setting := configMod.GetAllSettings(companyId).Data


    for key, ledger := range ledgerByName {
        collectionFilter := bson.M {
            "CompanyId": companyId,
            "LedgerName": key,
        }

        dbFilter := handlers.DocumentFilter {
            Ctx: context.TODO(),
            Filter: collectionFilter,
            UsePagination: false,
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
        billResponse := osMod.GetOutstandingCollection().FindDocuments(dbFilter)
        if billResponse.Err != nil {
            continue
        }
        var bills []osMod.Bill
        var totalAmount float64 = 0
        istLocation, _ := time.LoadLocation("Asia/Kolkata")
        for _, item := range billResponse.Data {
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
            days := int32(diff.Hours() / 24)

            var bill osMod.Bill = osMod.Bill {
                LedgerName: item["LedgerName"].(string),
                LedgerGroupName: item["LedgerGroupName"].(string),
                BillName: item["Name"].(string),
                DueDate: dueDate,
                BillDate: billDate,
                DelayDays: days,
            }
            var amount = utils.ParseFloat64(item["Amount"])
            bill.Amount = amount
            totalAmount += amount
            bills = append(bills, bill)
        }

        fmt.Printf("Some Address value: %v", ledger)

        content := ReminderBody {
            PartyName : key,
            Address: "",
            TotalAmount: totalAmount,
            Bills: bills,
        }
        emailBody := handlers.WriteToTemplate("/home/jayesh/development/research/templateWriter/osTemplate.html", content)


        // create email
        to := make([]string, 1)
        to[0] = "softgen.aquib.shaikh@gmail.com"

        cc := make([]string, 3)
        cc[0] = "rajeshri@4qs.in"
        cc[1] = "softgen.mustafa.khan@gmail.com"
        cc[2] = "softgen.saish.jagtap@gmail.com"

        var emailSettings = models.EmailSettings {
            To: to,
            Cc: cc,
            SmtpPort:"587",
            SmtpServer: "smtp.gmail.com",
            Subject: "Sample Email",
            Body: "Here is a sample email" + emailBody,
            BodyType: 1,
        }

        // send email 
        err := handlers.SendEmail(emailSettings)
        if err != nil  {
            fmt.Println("Failed to send email:", err)
        }
    }

}
