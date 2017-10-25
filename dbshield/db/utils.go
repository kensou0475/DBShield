package db

import (
	"errors"
	"reflect"
)

//Contains check if target contains obj
func Contains(target interface{}, obj interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in array")
}

//SContains string array contains
func SContains(target []string, obj string) (bool, error) {
	for _, t := range target {
		if t == obj {
			return true, nil
		}
	}
	return false, errors.New("not in array")
}
