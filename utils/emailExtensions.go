package utils
import (
    "strings"
    "reflect"
    "fmt"
)

func GenerateHTMLTable(data interface{}) (string, error) {
    val := reflect.ValueOf(data)

    // Ensure the input is a slice
    if val.Kind() != reflect.Slice {
        return "", fmt.Errorf("input data must be a slice")
    }

    // Ensure the slice elements are structs
    if val.Len() == 0 || val.Index(0).Kind() != reflect.Struct {
        return "", fmt.Errorf("slice elements must be structs")
    }

    var sb strings.Builder
    sb.WriteString("<table border='1'>\n")

    // Generate the table header based on struct field names
    sb.WriteString("<tr>")
    elemType := val.Index(0).Type()
    for i := 0; i < elemType.NumField(); i++ {
        fieldName := elemType.Field(i).Name
        sb.WriteString(fmt.Sprintf("<th>%s</th>", fieldName))
    }
    sb.WriteString("</tr>\n")

    // Generate the table rows based on struct field values
    for i := 0; i < val.Len(); i++ {
        sb.WriteString("<tr>")
        structVal := val.Index(i)
        for j := 0; j < structVal.NumField(); j++ {
            fieldVal := structVal.Field(j)
            sb.WriteString(fmt.Sprintf("<td>%v</td>", fieldVal.Interface()))
        }
        sb.WriteString("</tr>\n")
    }

    sb.WriteString("</table>")
    return sb.String(), nil
}
