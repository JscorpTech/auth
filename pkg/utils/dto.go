package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func getJSONFieldName(structField reflect.StructField) string {
	tag := structField.Tag.Get("json")
	if tag == "" {
		return structField.Name
	}
	return strings.Split(tag, ",")[0]
}

func FormatValidationErrors(err error, obj interface{}) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		val := reflect.TypeOf(obj).Elem()
		for _, fieldErr := range validationErrors {
			structField, _ := val.FieldByName(fieldErr.StructField())
			jsonName := getJSONFieldName(structField)

			switch fieldErr.Tag() {
			case "required":
				errorsMap[jsonName] = jsonName + " is required"
			case "email":
				errorsMap[jsonName] = "Invalid email format"
			case "min":
				errorsMap[jsonName] = jsonName + " is too short"
			default:
				errorsMap[jsonName] = "Invalid value"
			}
		}
	}

	return errorsMap
}
