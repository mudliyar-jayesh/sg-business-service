package utils

import (
	"encoding/json"
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
    UserId uint64
}

type ResponseStruct struct {
    Data any
    Count int
}

func NewResponseStruct[T any](data T, length int) ResponseStruct{
    return ResponseStruct {
        Data: data,
        Count: length,
    }
}

func (r ResponseStruct) ToJson(res http.ResponseWriter) error {
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    return json.NewEncoder(res).Encode(r)
}
