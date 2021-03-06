package transformer

import (
	"reflect"
	"strings"
)

const structTag = "keep"

type emptyStruct struct {
}

func contains(array []string, test string) bool {
	for _, i := range array {
		if i == test {
			return true
		}
	}
	return false
}

func getBsonName(t reflect.StructField) string {
	if bson := t.Tag.Get("bson"); len(bson) > 0 {
		bsonS := strings.Split(bson, ",")
		if len(bsonS) > 0 && len(bsonS[0]) > 0 {
			return bsonS[0]
		}
	}
	return t.Name
}

func getValue(v reflect.Value, field string, omitempty bool, t reflect.StructField, level int) interface{} {

	if t.Type.Name() == "ObjectID" {
		if omitempty && v.Len() == 0 {
			return nil
		}
		for i := v.Len() - 1; i >= 0; i-- {
			oidVal := v.Index(i).Uint()
			if oidVal != 0 {
				return v.Interface()
			}
			if oidVal == 0 && i == 0 {
				return nil
			}
		}
		return v
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		var array []interface{}

		for i := 0; i < v.Len(); i++ {
			cleaned := cleaner(v.Index(i).Interface(), field, level+1)
			array = append(array, cleaned)
		}

		if omitempty && len(array) == 0 {
			return nil
		}

		return array
	case reflect.Struct:
		return cleaner(v.Interface(), field, level+1)
	case reflect.String:
		if omitempty && v.Len() == 0 {
			return nil
		}
		return v.Interface()
	default:
		return v.Interface()
	}
}

func Clean(obj interface{}, field string) interface{} {
	return cleaner(obj, field, 1)
}

func cleaner(obj interface{}, field string, level int) interface{} {
	val, typ := reflect.ValueOf(obj), reflect.TypeOf(obj)
	returnValue := make(map[string]interface{})

	if val.Kind() != reflect.Struct {
		return obj
	}

	for i := 0; i < val.NumField(); i++ {
		v, t := val.Field(i), typ.Field(i)
		tags := strings.Split(t.Tag.Get(structTag), ",")

		if contains(tags, field) {
			key := getBsonName(t)
			omitempty := contains(tags, "omitempty")
			value := getValue(v, field, omitempty, t, level)

			if value != nil {
				returnValue[key] = value
			}
		}
	}

	if len(returnValue) == 0 {
		if level == 1 {
			return emptyStruct{}
		}
		return nil
	}
	return returnValue
}
