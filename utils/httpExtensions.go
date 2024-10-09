package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func GetValueBySortOrder(sortOrder string) int {
	if len(sortOrder) > 0 && sortOrder == "desc" {
		return -1
	}
	return 1
}

func GetBoolFromQuery(req *http.Request, queryParameter string) bool {
	boolStr := req.URL.Query().Get(queryParameter)
	if boolStr == "" {
		boolStr = "true"
	}
	var value bool
	value, _ = strconv.ParseBool(boolStr)
	return value
}

func ResolveHeaders(headers *http.Header) (RequestHeader, error) {
	var companyid string = headers.Get("companyid")
	var userid_str string = headers.Get("userid")

	userid, err := strconv.ParseUint(userid_str, 10, 64)

	return RequestHeader{CompanyId: companyid, UserId: userid}, err
}

func (headers *RequestHeader) HandleErrorOrIllegalValues(res http.ResponseWriter, err *error) bool {
	/**
	  Handles header errors and detects false headers like empty companyid or zero userid,
	  if found to be false header then writes to http.ResponseWriter and returns true
	  otherwise returns false.
	*/
	// TODO: Add additional checks for user-company_id linkages like BMRM middleware
	if *err != nil || len(headers.CompanyId) == 0 || headers.UserId == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Attempt to unauthorized access without secure headers"))
		return true
	}
	return false
}

func ReadRequestBody[T any](req *http.Request) (*T, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	var data T
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

type RequestHeader struct {
	CompanyId string
	UserId    uint64
}

type ResponseStruct struct {
	Data  any
	Count int
}

func NewResponseStruct[T any](data T, length int) ResponseStruct {
	return ResponseStruct{
		Data:  data,
		Count: length,
	}
}

func (r ResponseStruct) ToJson(res http.ResponseWriter) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	return json.NewEncoder(res).Encode(r)
}

// GetFromUms is a generic function to get data from a URL and unmarshal it into a slice of type T
func GetFromPortal[T any](url string, headers RequestHeader) *T {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil
	}

	// Set custom headers
	req.Header.Set("companyid", headers.CompanyId)

	userId := strconv.FormatUint(headers.UserId, 10)
	req.Header.Set("userid", userId) // Assuming UserId is already a string

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return nil
	}
	defer resp.Body.Close() // Ensure response body is closed after reading

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil
	}

	// Unmarshal the response body into a slice of type T
	var data T
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil
	}

	return &data
}
