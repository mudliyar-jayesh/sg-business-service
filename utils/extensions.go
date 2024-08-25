package utils

import (
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
)

func GroupBy(documents []bson.M, key string) map[string][]bson.M {
    grouped := make(map[string][]bson.M)

    for _, doc := range documents {
        field, exists := doc[key].(string)
        if !exists {
            fmt.Println("LedgerName field missing or not a string")
            continue
        }
        grouped[field] = append(grouped[field], doc)
    }

    return grouped
}

func Intersection(slice1, slice2 []string) []string {
	// Create a map to store elements from the first slice
	elementMap := make(map[string]struct{}, len(slice1))

	// Populate the map with elements from the first slice
	for _, item := range slice1 {
		elementMap[item] = struct{}{}
	}

	// Create a slice to store the intersection results
	var result []string

	// Check elements of the second slice against the map
	for _, item := range slice2 {
		if _, found := elementMap[item]; found {
			result = append(result, item)
		}
	}

	return result
}
