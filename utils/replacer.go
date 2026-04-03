package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var templateRegex = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.]+)\s*\}\}`)

func RenderTemplate(text string, data any) string {
	return templateRegex.ReplaceAllStringFunc(text, func(match string) string {
		submatch := templateRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}
		keyPath := submatch[1]

		val, found := getValueByPath(keyPath, data)
		if !found || val == "" {
			return match
		}

		return fmt.Sprintf("%v", val)
	})
}

func getValueByPath(path string, data any) (any, bool) {
	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys {
		v := reflect.ValueOf(current)

		for v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return nil, false
			}
			v = v.Elem()
		}

		switch v.Kind() {

		case reflect.Map:

			mapVal := v.MapIndex(reflect.ValueOf(key))
			if !mapVal.IsValid() {
				return nil, false
			}
			current = mapVal.Interface()

		case reflect.Struct:
			var field reflect.Value
			foundField := false

			for i := 0; i < v.NumField(); i++ {
				typeField := v.Type().Field(i)
				if strings.EqualFold(typeField.Name, key) {
					field = v.Field(i)
					foundField = true
					break
				}
			}

			if !foundField {
				return nil, false
			}

			if !field.CanInterface() {
				return nil, false
			}
			current = field.Interface()

		default:
			return nil, false
		}
	}

	return current, true
}
