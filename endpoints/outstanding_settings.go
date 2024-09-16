package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"sg-business-service/handlers"
	"sg-business-service/models"
	"sg-business-service/modules/outstanding"
	"sg-business-service/modules/outstanding/reminders"
	"sg-business-service/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateOsSetting(res http.ResponseWriter, req *http.Request) {
	collection := handlers.GetCollection("BMRM", "OutstandingSettings")
	var mongoHandler = handlers.NewMongoHandler(collection)

	companyId := req.Header.Get("CompanyId")

	body, err := utils.ReadRequestBody[outstanding.OsShareSettings](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

	body.CompanyId = companyId

	docFilter := handlers.DocumentFilter{
		Ctx: context.TODO(),
		Filter: bson.M{
			"CompanyId": companyId,
		},
		UsePagination: false,
		Limit:         0,
		Offset:        0,
	}

	var results handlers.DocumentResponse = mongoHandler.FindDocuments(docFilter)
	fmt.Println(results.Data)
	if len(results.Data) != 0 {
		http.Error(res, "Already Exists", http.StatusBadRequest)
		return
	}

	_, err = handlers.InsertDocument("BMRM", "OutstandingSettings", body)
	if err != nil {
		http.Error(res, "Could not create entry", http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func UpdateOsSetting(res http.ResponseWriter, req *http.Request) {
	collection := handlers.GetCollection("BMRM", "OutstandingSettings")
	var mongoHandler = handlers.NewMongoHandler(collection)

	companyId := req.Header.Get("CompanyId")

	body, err := utils.ReadRequestBody[outstanding.OsShareSettings](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

	body.CompanyId = companyId

	filter := bson.M{
		"CompanyId": companyId,
	}

	update := bson.M{
		"$set": body,
	}

	mongoHandler.UpdateDocument("BMRM", "OutstandingSettings", filter, update)
	if err != nil {
		http.Error(res, "Could not create entry", http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func GetSetting(res http.ResponseWriter, req *http.Request) {
	collection := handlers.GetCollection("BMRM", "OutstandingSettings")
	var mongoHandler = handlers.NewMongoHandler(collection)

	companyId := req.Header.Get("CompanyId")

	docFilter := handlers.DocumentFilter{
		Ctx: context.TODO(),
		Filter: bson.M{
			"CompanyId": companyId,
		},
		UsePagination: false,
		Limit:         0,
		Offset:        0,
	}

	var results handlers.DocumentResponse = mongoHandler.FindDocuments(docFilter)
	if results.Err != nil {
		http.Error(res, "No Data", http.StatusBadRequest)
		return
	}
	response := utils.NewResponseStruct(results.Data, len(results.Data))
	response.ToJson(res)

}

func SendEmail(res http.ResponseWriter, req *http.Request) {

	to := make([]string, 1)
	to[0] = "softgen.aquib.shaikh@gmail.com"
	cc := make([]string, 1)
	cc[0] = "jayeshmudlyiar2112000@gmail.com"

	var emailSettings = models.EmailSettings{
		To:         to,
		Cc:         cc,
		SmtpPort:   "587",
		SmtpServer: "smtp.gmail.com",
		Subject:    "Sample Email",
		Body:       "Here is a sample email",
		BodyType:   1,
	}
	err := handlers.SendEmail(emailSettings)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		http.Error(res, "Could not send email", http.StatusBadRequest)
	}
}

func SendLedgerEmail(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")

	parties, err := utils.ReadRequestBody[EmailRequest](req)

	if err != nil {
		fmt.Println("Error parsing the list of ledgerNames")
	}

	reminders.SendEmailReminder(companyId, parties.Parties)
}

type EmailRequest struct {
	Parties []string
}
