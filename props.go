package ctl

import (
	"reflect"
	"strings"
)

func setProperty(name string, value string, obj interface{}) bool {
	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Ptr {
		return false
	}
	structVal := val.Elem()
	for i := 0; i < structVal.NumField(); i++ {
		valueField := structVal.Field(i)
		typeField := structVal.Type().Field(i)
		if strings.ToLower(typeField.Name) == strings.ToLower(name) {
			if valueField.IsValid() && valueField.CanSet() && valueField.Kind() == reflect.String {
				valueField.SetString(value)
				return true
			}
		}
	}
	return false
}
