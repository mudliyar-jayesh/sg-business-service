package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"strconv"
	"strings"
	"unicode"
)

// FloatFromString represents a float64 value that may come as a string in BSON.
type FloatFromString struct {
	Value float64
}

func cleanString(s string) string {
	// Remove non-printable characters and trim whitespace
	var cleaned strings.Builder
	for _, r := range s {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			cleaned.WriteRune(r)
		}
	}
	return strings.TrimSpace(cleaned.String())
}

// UnmarshalBSONValue is a custom BSON unmarshaler for FloatFromString.
func (f *FloatFromString) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var err error

	var strValue string = string(data)
	cleanedStr := cleanString(strValue)
	floatValue, err := strconv.ParseFloat(cleanedStr, 64)

	if err != nil {
		fmt.Println("Error converting string to float64:", err)
		return err
	}
	f.Value = floatValue
	return nil
}
