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
