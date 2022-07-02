package goval

import (
	"reflect"
)

//Set Update the structure field with a given value.
func Set[T any](target any, path Path, newValue T) {
	SetFunc(target, path, func(_ T, _ PathInfo) T {
		return newValue
	})
}

// SetFunc Update the structure field with a function value.
func SetFunc[T any](target any, path Path, fn func(v T, pathInfo PathInfo) T) {
	refTarget := reflect.ValueOf(target) // *interface{}
	pathInfo := PathInfo{
		RequirePath: path,
	}
	each(refTarget, path.Split(), pathInfo, func(v any, pathInfo PathInfo) {
		newVal := fn(v.(T), pathInfo)
		pathInfo.fieldValue.Set(reflect.ValueOf(newVal))
	})
}
