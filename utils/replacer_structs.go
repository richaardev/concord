package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func FillStructPlaceholders[T any](input T, placeholders map[string]interface{}) (T, error) {
	data := nestedMap(placeholders)

	if s, ok := any(input).(string); ok {
		rendered, err := renderTemplate(s, data)
		if err != nil {
			return input, err
		}

		result, _ := any(rendered).(T)
		return result, nil
	}

	val := reflect.ValueOf(input)
	filled, err := fillValue(val, data)
	if err != nil {
		return input, err
	}

	result, ok := filled.Interface().(T)
	if !ok {
		return input, fmt.Errorf("tipo inesperado após preenchimento: %T", filled.Interface())
	}
	return result, nil
}

func fillValue(v reflect.Value, data map[string]interface{}) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return v, nil
		}
		inner, err := fillValue(v.Elem(), data)
		if err != nil {
			return v, err
		}
		ptr := reflect.New(inner.Type())
		ptr.Elem().Set(inner)
		return ptr, nil

	case reflect.Interface:
		if v.IsNil() {
			return v, nil
		}
		inner, err := fillValue(v.Elem(), data)
		if err != nil {
			return v, err
		}

		iface := reflect.New(v.Type()).Elem()
		if inner.Type().Implements(v.Type()) {
			iface.Set(inner)
		} else if inner.Type().AssignableTo(v.Type()) {
			iface.Set(inner)
		} else {
			return v, nil
		}
		return iface, nil

	case reflect.String:
		rendered, err := renderTemplate(v.String(), data)
		if err != nil {
			return v, err
		}
		return reflect.ValueOf(rendered).Convert(v.Type()), nil

	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		out.Set(v)
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			filled, err := fillValue(v.Field(i), data)
			if err != nil {
				return v, fmt.Errorf("campo %s: %w", field.Name, err)
			}
			if filled.Type().AssignableTo(out.Field(i).Type()) {
				out.Field(i).Set(filled)
			}
		}
		return out, nil

	case reflect.Slice:
		if v.IsNil() {
			return v, nil
		}
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			filled, err := fillValue(v.Index(i), data)
			if err != nil {
				return v, fmt.Errorf("índice %d: %w", i, err)
			}
			if filled.Type().AssignableTo(out.Index(i).Type()) {
				out.Index(i).Set(filled)
			} else {
				out.Index(i).Set(v.Index(i))
			}
		}
		return out, nil

	case reflect.Array:
		out := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			filled, err := fillValue(v.Index(i), data)
			if err != nil {
				return v, fmt.Errorf("índice %d: %w", i, err)
			}
			if filled.Type().AssignableTo(out.Index(i).Type()) {
				out.Index(i).Set(filled)
			}
		}
		return out, nil

	case reflect.Map:
		if v.IsNil() {
			return v, nil
		}
		out := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			filled, err := fillValue(v.MapIndex(key), data)
			if err != nil {
				return v, fmt.Errorf("chave %v: %w", key, err)
			}
			if filled.Type().AssignableTo(v.Type().Elem()) {
				out.SetMapIndex(key, filled)
			} else {
				out.SetMapIndex(key, v.MapIndex(key))
			}
		}
		return out, nil

	default:

		return v, nil
	}
}

func renderTemplate(tmplStr string, data map[string]interface{}) (string, error) {
	if !strings.Contains(tmplStr, "{{") {
		return tmplStr, nil
	}

	// funcMap := template.FuncMap{
	// 	"default": func(def, val interface{}) interface{} {
	// 		if val == nil || val == "" {
	// 			return def
	// 		}
	// 		return val
	// 	},
	// 	"upper": strings.ToUpper,
	// 	"lower": strings.ToLower,
	// }
	//
	// tmpl, err := template.New("fill").
	// 	Funcs(funcMap).
	// 	Option("missingkey=zero").
	// 	Parse(tmplStr)
	// if err != nil {
	// 	return "", fmt.Errorf("template parse error: %w", err)
	// }
	//
	// var buf bytes.Buffer
	// if err := tmpl.Execute(&buf, data); err != nil {
	// 	return "", fmt.Errorf("template execute error: %w", err)
	// }

	result := RenderTemplate(tmplStr, data)
	return result, nil
}

func nestedMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = deepConvert(v)
	}
	return result
}

func deepConvert(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return nil
		}
		return deepConvert(rv.Elem().Interface())

	case reflect.Map:
		out := make(map[string]interface{})
		for _, k := range rv.MapKeys() {
			out[fmt.Sprintf("%v", k.Interface())] = deepConvert(rv.MapIndex(k).Interface())
		}
		return out

	case reflect.Slice, reflect.Array:
		out := make([]interface{}, rv.Len())
		for i := range out {
			out[i] = deepConvert(rv.Index(i).Interface())
		}
		return out

	case reflect.Struct:
		out := make(map[string]interface{})
		t := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			if !t.Field(i).IsExported() {
				continue
			}
			out[t.Field(i).Name] = deepConvert(rv.Field(i).Interface())
		}
		return out

	default:
		return v
	}
}
