package utils

import (
    "go.mongodb.org/mongo-driver/bson"
    "regexp"
    //"strings"
)

func GenerateSearchFilter(searchText, searchKey string) []bson.M{
	// Regular expression to match numbers, single letters, and words
	re := regexp.MustCompile(`\d+|[a-zA-Z]+`)

	// Find all matching tokens
	tokens := re.FindAllString(searchText, -1)
    //tokens := strings.Fields(searchText)

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
