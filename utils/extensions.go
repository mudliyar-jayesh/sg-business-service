package utils

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
	"reflect"
	"sort"
	"strconv"
	"time"
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
func ToLookup[T any, K comparable](items []T, keySelector func(T) K) map[K][]T {
	dict := make(map[K][]T)
	for _, item := range items {
		key := keySelector(item)
		var values = dict[key]
		values = append(values, item)
		dict[key] = values
	}
	return dict
}

func ToDict[T any, K comparable](items []T, keySelector func(T) K) map[K]T {
	dict := make(map[K]T)
	for _, item := range items {
		key := keySelector(item)
		dict[key] = item
	}
	return dict
}
func Where[S any](source []S, predicate func(S) bool) []S {
	var result []S
	for _, s := range source {
		if predicate(s) {
			result = append(result, s)
		}
	}
	return result
}

func Select[S any, T any](source []S, selector func(S) T) []T {
	result := make([]T, len(source))
	for i, s := range source {
		result[i] = selector(s)
	}
	return result
}

func GroupFor[S any, K comparable](source []S, keySelector func(S) K) map[K][]S {
	result := make(map[K][]S)
	for _, s := range source {
		key := keySelector(s)
		result[key] = append(result[key], s)
	}
	return result
}

func GroupBySelect[S any, K comparable, V any](source []S, keySelector func(S) K, elementSelector func(S) V) map[K][]V {
	result := make(map[K][]V)
	for _, s := range source {
		key := keySelector(s)
		value := elementSelector(s)
		result[key] = append(result[key], value)
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
func ProcessBatch[T any](values []T, chunkSize int, predicate func([]T)) {
	var defaultChunkSize = 100
	var dataLength = len(values)

	var workingChuckSize int = chunkSize
	if dataLength < chunkSize {
		workingChuckSize = defaultChunkSize
		if dataLength < defaultChunkSize {
			workingChuckSize = dataLength
		}
	}

	var chunkCount int = int(math.Ceil(float64(dataLength / workingChuckSize)))

	for chunkNumber := 0; chunkNumber < chunkCount; chunkNumber++ {
		var slice []T = getChunk(values, workingChuckSize, chunkNumber)
		predicate(slice)
	}

}

func getChunk[T any](list []T, chunkSize, chunkNumber int) []T {
	startIndex := chunkNumber * chunkSize
	if startIndex >= len(list) {
		return nil
	}

	endIndex := startIndex + chunkSize
	if endIndex > len(list) {
		endIndex = len(list)
	}

	return list[startIndex:endIndex]
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ToHashSet converts a slice of type T to a set (implemented as a map[T]struct{})
// It ensures that all values in the returned map are unique.
func ToHashSet[T comparable](items []T) map[T]struct{} {
	set := make(map[T]struct{})
	for _, item := range items {
		set[item] = struct{}{} // Using struct{}{} to save memory
	}
	return set
}
func Distinct[T comparable](items []T) []T {
	set := make(map[T]struct{})
	var result []T
	for _, item := range items {
		if _, exists := set[item]; !exists {
			set[item] = struct{}{}        // Mark the item as seen
			result = append(result, item) // Add the distinct item to the result slice
		}
	}
	return result
}
func ContainsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func SortBy[T any](slice []T, lessFunc func(i, j T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return lessFunc(slice[i], slice[j])
	})
}

func Paginate[T any](slice []T, limit, offset int) []T {
	// Check bounds to avoid slicing out of range
	if offset > len(slice) {
		return []T{} // Return empty slice if offset exceeds slice length
	}

	end := offset + limit

	// Ensure the end index does not exceed the slice length
	if end > len(slice) {
		end = len(slice)
	}

	// Return the sub-slice
	return slice[offset:end]
}

func SortByField(slice interface{}, fieldName string, sortByAsc bool) {
	// Get the value of the slice
	v := reflect.ValueOf(slice)

	// Check if the passed interface is a slice
	if v.Kind() != reflect.Slice {
		log.Fatalf("SortByField: expected a slice, got %T", slice)
	}

	// Get the element type of the slice
	elemType := v.Type().Elem()

	// Sort the slice using the sort.Slice function
	sort.Slice(slice, func(i, j int) bool {
		// Get the i-th and j-th elements of the slice
		vi := v.Index(i)
		vj := v.Index(j)

		// Handle if elements are pointers to structs
		if elemType.Kind() == reflect.Ptr {
			vi = vi.Elem() // Dereference pointer to access struct
			vj = vj.Elem()
		}

		// Ensure the elements of the slice are structs
		if vi.Kind() != reflect.Struct || vj.Kind() != reflect.Struct {
			log.Fatalf("SortByField: expected a slice of structs or pointers to structs, got %s", elemType.Kind())
		}

		// Get the values of the specified field for the i-th and j-th elements
		fieldI := vi.FieldByName(fieldName)
		fieldJ := vj.FieldByName(fieldName)

		// Check if the field exists
		if !fieldI.IsValid() || !fieldJ.IsValid() {
			log.Fatalf("SortByField: field %s not found in struct %s", fieldName, elemType.Name())
		}

		// Handle pointer fields, considering nil values
		if fieldI.Kind() == reflect.Ptr {
			// Handle nil cases: if one is nil and the other is not, prioritize non-nil value
			if fieldI.IsNil() && !fieldJ.IsNil() {
				return sortByAsc
			}
			if !fieldI.IsNil() && fieldJ.IsNil() {
				return !sortByAsc
			}
			// If both are non-nil, dereference the pointers
			if !fieldI.IsNil() && !fieldJ.IsNil() {
				fieldI = fieldI.Elem()
				fieldJ = fieldJ.Elem()
			}
		}

		// Compare the field values based on their kind
		var result bool
		switch fieldI.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result = fieldI.Int() < fieldJ.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// Compare unsigned integers, considering nil as the smallest (0)
			result = fieldI.Uint() < fieldJ.Uint()
		case reflect.Float32, reflect.Float64:
			result = fieldI.Float() < fieldJ.Float()
		case reflect.String:
			result = fieldI.String() < fieldJ.String()
		case reflect.Bool:
			result = fieldI.Bool() && !fieldJ.Bool()
		case reflect.Struct:
			// Handle time.Time comparison
			if fieldI.Type() == reflect.TypeOf(time.Time{}) {
				timeI := fieldI.Interface().(time.Time)
				timeJ := fieldJ.Interface().(time.Time)
				result = timeI.Before(timeJ) // Sort in ascending order (oldest first)
			} else {
				log.Fatalf("SortByField: unsupported struct type %s", fieldI.Type())
			}
		default:
			log.Fatalf("SortByField: unsupported field type %s", fieldI.Kind())
		}

		// If sortByAsc is false, reverse the result for descending order
		if !sortByAsc {
			return !result
		}
		return result
	})
}

type DurationKeys struct {
	StartDate time.Time
	EndDate   time.Time
	Keys      []string
}

// GetDurationKeys generates duration keys (Daily, Weekly, Monthly, Yearly) with pagination
func GetDurationKeys(durationType string, pageSize, page int) DurationKeys {
	var keys []string
	today := time.Now()
	startDate := today
	endDate := today
	location, _ := time.LoadLocation("Asia/Kolkata") // Use IST location

	switch durationType {
	case "Daily":
		keys, startDate, endDate = generateDailyKeys(today, location, pageSize, page)
	case "Weekly":
		keys, startDate, endDate = generateWeeklyKeys(today, location, pageSize, page)
	case "Monthly":
		keys, startDate, endDate = generateMonthlyKeys(today, location, pageSize, page)
	case "Yearly":
		keys, startDate, endDate = generateYearlyKeys(today, location, pageSize, page)
	default:
		fmt.Println("Unknown duration type")
	}

	return DurationKeys{
		StartDate: startDate,
		EndDate:   endDate,
		Keys:      keys,
	}
}

func generateDailyKeys(today time.Time, location *time.Location, pageSize, page int) ([]string, time.Time, time.Time) {
	keys := []string{}
	for i := 0; i < pageSize; i++ {
		day := today.AddDate(0, 0, -((page-1)*pageSize + i))
		keys = append(keys, day.Format("02-Jan-2006"))
	}
	startDate := today.AddDate(0, 0, -((page - 1) * pageSize))
	endDate := today.AddDate(0, 0, -(page * pageSize))
	return keys, endDate, startDate
}

func generateWeeklyKeys(today time.Time, location *time.Location, pageSize, page int) ([]string, time.Time, time.Time) {
	keys := []string{}
	for i := 0; i < pageSize; i++ {
		weekOffset := (page-1)*pageSize + i
		startOfWeek := today.AddDate(0, 0, -int(today.Weekday())-7*weekOffset+1)
		endOfWeek := startOfWeek.AddDate(0, 0, 6)
		key := fmt.Sprintf("%s to %s", startOfWeek.Format("2 Jan 2006"), endOfWeek.Format("2 Jan 2006"))
		keys = append(keys, key)
	}
	startDate := today.AddDate(0, 0, -(page-1)*7*pageSize)
	endDate := today.AddDate(0, 0, -page*7*pageSize)
	return keys, endDate, startDate
}

func generateMonthlyKeys(today time.Time, location *time.Location, pageSize, page int) ([]string, time.Time, time.Time) {
	keys := []string{}
	for i := 0; i < pageSize; i++ {
		month := today.AddDate(0, -(page-1)*pageSize-i, 0)
		key := month.Format("Jan-2006")
		keys = append(keys, key)
	}
	startDate := today.AddDate(0, -(page-1)*pageSize, 0)
	endDate := today.AddDate(0, -page*pageSize, 0)
	return keys, endDate, startDate
}

func generateYearlyKeys(today time.Time, location *time.Location, pageSize, page int) ([]string, time.Time, time.Time) {
	keys := []string{}
	for i := 0; i < pageSize; i++ {
		year := today.AddDate(-(page-1)*pageSize-i, 0, 0)
		key := year.Format("2006")
		keys = append(keys, key)
	}
	startDate := today.AddDate(-(page-1)*pageSize, 0, 0)
	endDate := today.AddDate(-page*pageSize, 0, 0)
	return keys, endDate, startDate
}
