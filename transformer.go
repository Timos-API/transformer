package transformer

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

const structTag = "keep"

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

func Clean(obj interface{}, key string) bson.M {
	val, typ, returnValue := reflect.ValueOf(obj), reflect.TypeOf(obj), bson.M{}

	for i := 0; i < val.NumField(); i++ {
		v, t := val.Field(i), typ.Field(i)
		tags := strings.Split(t.Tag.Get(structTag), ",")

		if contains(tags, key) {
			name := getBsonName(t)

			if v.Kind() == reflect.Struct {
				returnValue[name] = Clean(v.Interface(), key)
			} else {
				returnValue[name] = v.Interface()
			}

		}
	}
	return returnValue
}
