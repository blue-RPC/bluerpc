package bluerpc

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func GoFieldsToTSObj(someStruct reflect.Type) string {
	stringBuilder := strings.Builder{}

	if someStruct.Kind() == reflect.Ptr {
		someStruct = someStruct.Elem()
	} else if someStruct.Kind() != reflect.Struct {
		return "any"
	}

	stringBuilder.WriteString("{")

	for i := 0; i < someStruct.NumField(); i++ {
		field := someStruct.Field(i)
		fieldName := field.Name
		fieldType := field.Type

		json := field.Tag.Get("json")
		input := field.Tag.Get("input")

		if input != "" {
			regex := regexp.MustCompile("[^a-zA-Z]+")
			fieldName = regex.ReplaceAllString(input, "")
		}

		if json != "" {
			regex := regexp.MustCompile("[^a-zA-Z]+")
			fieldName = regex.ReplaceAllString(json, "")
		}

		// Append TypeScript field definition to the StringBuilder
		stringBuilder.WriteString(fmt.Sprintf(" %s: %s", fieldName, goTypeToTSType(fieldType)))
		tags := field.Tag.Get("validate")
		if !strings.Contains(tags, "required") {
			stringBuilder.WriteString("|undefined")
		}
		stringBuilder.WriteString(",")

	}
	stringBuilder.WriteString("}")
	return stringBuilder.String()
}

func goTypeToTSType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		elemType := goTypeToTSType(t.Elem())
		return "Array<" + elemType + ">"
		// Add more type mappings as needed
	case reflect.Map:
		keyType := goTypeToTSType(t.Key())
		valueType := goTypeToTSType(t.Elem())
		return fmt.Sprintf("Record<%s, %s>", keyType, valueType)
	default:
		return "any"
	}
}
