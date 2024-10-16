package reminders

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sg-business-service/config"
	"sg-business-service/handlers"
	"sg-business-service/models"
	ledgersMod "sg-business-service/modules/ledgers"
	osMod "sg-business-service/modules/outstanding"
	configMod "sg-business-service/modules/outstanding/settings"
	"sg-business-service/utils"
	"time"

	"github.com/leekchan/accounting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ac = accounting.Accounting{Symbol: "₹", Precision: 2}

func SendEmailReminder(companyId string, ledgerNames []string) error {
	ledgers := ledgersMod.GetByNames(companyId, ledgerNames)
	ledgerByName := utils.ToDict(ledgers, func(ledger ledgersMod.MetaLedger) string {
		return ledger.Name
	})

	// Get Outstanding settings from Database for this company
	settingData, settingErr := configMod.GetAllSettings(companyId)
	if settingData == nil || len(settingData) > 1 {
		return settingErr
	}
	setting := settingData[0]

	for key, _ := range ledgerByName {
		collectionFilter := bson.M{
			"CompanyId":  companyId,
			"LedgerName": key,
		}

		// if cutoff date is not present then apply billDate filter.
		if len(setting.CutOffDate) > 0 {
			parsedDate, err := time.Parse("2006-01-02", setting.CutOffDate)

			if err != nil {
				fmt.Println("Error occoured while parsing date", setting.CutOffDate)
			}

			collectionFilter = bson.M{
				"CompanyId":  companyId,
				"LedgerName": key,
				"BillDate.Date": bson.M{
					"$gte": parsedDate,
				},
			}
		}

		dbFilter := handlers.DocumentFilter{
			Ctx:           context.TODO(),
			Filter:        collectionFilter,
			UsePagination: false,
			Projection: bson.M{
				"LedgerName":      1,
				"LedgerGroupName": 1,
				"BillDate":        "$BillDate.Date",
				"DueDate":         "$BillCreditPeriod.DueDate",
				"Amount":          "$ClosingBal.Amount",
				"OpeningAmount":   "$OpeningBal.Amount",
				"Name":            "$Name",
				"_id":             0,
			},
		}
		billResponse := osMod.GetOutstandingCollection().FindDocuments(dbFilter)
		if billResponse.Err != nil {
			continue
		}
		var bills []osMod.Bill
		var totalAmount float64 = 0
		istLocation, _ := time.LoadLocation("Asia/Kolkata")

		layout := "02 January 2006"

		for _, item := range billResponse.Data {
			billDateValue := item["BillDate"].(primitive.DateTime).Time()
			billDate := billDateValue.In(istLocation).Format(layout)

			var dueDate string
			dueDateValue := item["DueDate"]
			if dueDateValue == nil {
				dueDate = billDate
			} else {
				dueDate = dueDateValue.(primitive.DateTime).Time().In(istLocation).Format(layout)
			}

			parsedTime, _ := time.Parse(layout, dueDate)
			today := time.Now().UTC()

			// Calculate the difference
			diff := today.Sub(parsedTime)

			// Get the difference in days
			days := int32(diff.Hours() / 24)

			// Additional filters
			eligible := false

			if (setting.SendAllDue) ||
				(setting.SendDueOnly && days < int32(setting.OverDueDays)) ||
				(setting.SendOverDueOnly && days >= int32(setting.OverDueDays)) {
				eligible = true
			}

			if !eligible {
				continue
			}

			var bill osMod.Bill = osMod.Bill{
				LedgerName:      item["LedgerName"].(string),
				LedgerGroupName: item["LedgerGroupName"].(string),
				BillName:        item["Name"].(string),
				DueDate:         dueDate,
				BillDate:        billDate,
				DelayDays:       days,
			}

			var amount = utils.ParseFloat64(item["Amount"])
			bill.Amount = amount
			bill.AmountStr = ac.FormatMoney(bill.Amount)
			totalAmount += amount
			bills = append(bills, bill)
		}

		content := ReminderBody{
			PartyName:      key,
			Address:        "",
			TotalAmount:    totalAmount,
			TotalAmountStr: ac.FormatMoney(totalAmount),
			Bills:          bills,
		}

		//emailBody := handlers.WriteToTemplate("C:\\Users\\softg\\Projects\\sg-business-service\\osTemplate.html", content)
		var templatePath string = config.LoadEmailTemplate()

		_, err := os.Stat(templatePath)

		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("\n%v path does not exists\n", templatePath)
			return err
		}

		var emailBody string
		if setting.TemplateName != nil || len(*setting.TemplateName) > 0 {
			template := GetByTemplateName(companyId, *setting.TemplateName)
			if template != nil && len(template.HtmlContent) > 0 {
				emailBody = handlers.WriteToTemplateFromString(template.HtmlContent, content)
			} else {
				emailBody = handlers.WriteToTemplate(templatePath, content)
			}
		} else {
			emailBody = handlers.WriteToTemplate(templatePath, content)
		}
		emailSetting := setting.EmailSetting

		var emailSettings = models.EmailSettings{
			To:         emailSetting.To,
			Cc:         emailSetting.Cc,
			SmtpPort:   emailSetting.SmtpPort,
			SmtpServer: emailSetting.SmtpServer,
			Subject:    emailSetting.Subject,
			Body:       emailSetting.Body + emailBody,
			BodyType:   1,
		}

		// send email
		err = handlers.SendEmail(emailSettings)
		if err != nil {
			fmt.Println("Failed to send email:", err)
			return err
		}
	}
	return nil
}
