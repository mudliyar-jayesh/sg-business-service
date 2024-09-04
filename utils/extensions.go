package utils

import (
    "strconv"
    "fmt"
    "errors"
    "reflect"
    "go.mongodb.org/mongo-driver/bson"
)

func GroupByKey[T any](items []T, key string) map[string][]T {
	grouped := make(map[string][]T)

	for _, item := range items {
		// Use reflection to get the field value
		val := reflect.ValueOf(item)
		field := val.FieldByName(key)

		// Ensure the field exists and is a string
		if !field.IsValid() || field.Kind() != reflect.String {
			fmt.Println(key, "field missing or not a string")
			continue
		}

		// Group items by the field value
		fieldValue := field.String()
		grouped[fieldValue] = append(grouped[fieldValue], item)
	}

	return grouped
}

// AggregateMultiFields performs multiple aggregations on the grouped data
func AggregateMultiFields[T any, R any](grouped map[string][]T, aggregationFuncs map[string]func([]T) R) map[string]map[string]R {
	results := make(map[string]map[string]R)

	for key, items := range grouped {
		results[key] = make(map[string]R)
		for aggKey, aggFunc := range aggregationFuncs {
			results[key][aggKey] = aggFunc(items)
		}
	}

	return results
}

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

func ToDictionary(bsonSlice []bson.M, keyField string) (map[string]interface{}, error) {
    dict := make(map[string]interface{})

    for _, item := range bsonSlice {
        key, ok := item[keyField].(string)
        if !ok {
            return nil, errors.New("invalid key string")
        }
        dict[key] = item
    }
    return dict, nil
}


func ParseFloat64(value interface{}) float64 {
    var result float64
    switch v := value.(type) {
    case float64:
        result = v
    case int:
        result = float64(v)
    case string:
        parsed, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return 0 // Return default value on error
        }
        result = parsed
    default:
        return 0 // Return default value if type is not handled
    }
    return result
}

