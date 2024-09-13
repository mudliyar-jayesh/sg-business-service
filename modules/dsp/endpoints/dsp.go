package endpoints

import (
	"encoding/base64"
	"encoding/csv"
	"net/http"
	"sg-business-service/models"
	"sg-business-service/modules/dsp"
	"sg-business-service/utils"
	"strings"
)

func GetStates(res http.ResponseWriter, req *http.Request) {
	states := dsp.GetStates()
	response := utils.NewResponseStruct(states, len(states))
	response.ToJson(res)
}

func UploadCsvFile(res http.ResponseWriter, req *http.Request) {
	body, err := utils.ReadRequestBody[models.File](req)
	if err != nil {
		http.Error(res, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(body.Data)
	if err != nil {
		http.Error(res, "Invalid File Data", http.StatusBadRequest)
		return
	}
	fileContent := string(fileBytes)

	reader := csv.NewReader(strings.NewReader(fileContent))
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(res, "Error Reading Csv File", http.StatusBadRequest)
		return
	}

	var entries []interface{}

	for index, record := range records {
		if index == 0 {
			continue
		}

		// Create a map (or slice) to represent the record instead of using a struct
		newEntry := map[string]string{
			"CircleName":   record[0],
			"RegionName":   record[1],
			"DivisionName": record[2],
			"OfficeName":   record[3],
			"Pincode":      record[4],
			"OfficeType":   record[5],
			"Delivery":     record[6],
			"District":     record[7],
			"StateName":    record[8],
			"Latitude":     record[9],
			"Longitude":    record[10],
		}

		// Append each entry to the entries slice
		entries = append(entries, newEntry)
	}

	// Now insert the entries slice using InsertMany or a custom batch method
	handler := dsp.GetCollection()
	handler.InsertMany(entries)
}
