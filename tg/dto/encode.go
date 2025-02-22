package dto

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func MarshalMap(d interface{}) (map[string]string, error) {
	resultMap := make(map[string]string)
	v := reflect.ValueOf(d)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		key, _, _ := strings.Cut(jsonTag, ",") // Handle JSON tag options
		if key == "" {
			key = field.Name // Use struct field name if no JSON tag
		}

		if isEncodeableStruct(fieldValue) {
			// Handle nested struct
			jsonData, err := json.Marshal(fieldValue.Interface())
			if err != nil {
				return nil, err
			}

			resultMap[key] = string(jsonData)
		} else {
			// Handle basic types
			resultMap[key] = fmt.Sprintf("%v", fieldValue.Interface())
		}
	}

	return resultMap, nil
}

// Helper function to check if a value is a struct that needs encoding
func isEncodeableStruct(v reflect.Value) bool {
	kind := v.Kind()

	// Dereference pointers
	if kind == reflect.Ptr {
		if v.IsNil() {
			return false
		}

		v = v.Elem()
		kind = v.Kind()
	}

	if kind != reflect.Struct {
		return false
	}

	return true
}
