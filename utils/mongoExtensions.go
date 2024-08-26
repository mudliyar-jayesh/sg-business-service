package utils

import (
    "go.mongodb.org/mongo-driver/bson"
    "strings"
)

func GenerateSearchFilter(searchText, searchKey string) []bson.M{
    tokens := strings.Fields(searchText)

    var regexFilters []bson.M

    for _, token := range tokens {
        regexFilters = append(regexFilters, bson.M {
            searchKey: bson.M {
                "$regex": token, 
                "$options": "i",// case insensitie
            },
        })
    }
    return regexFilters
}
