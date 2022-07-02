package goval

import (
	"reflect"
)

// funcEach a callback func given to each function.
//
// v: Value of the field specified in the path.
// owner: Structure owning the value pointer.
type funcEach func(v any, pathInfo PathInfo)

// Each executes the given function once for each field specified in the path.
func Each(target any, path Path, fn funcEach) {
	refTarget := reflect.ValueOf(target)
	pathInfo := PathInfo{
		RequirePath: path,
	}
	each(refTarget, path.Split(), pathInfo, fn)
}

// each executes the given function once for each field specified in the path.
//
// refTargetAddr: 参照のValueである必要がある
func each(target reflect.Value, paths []Path, pathInfo PathInfo, fn funcEach) {
	if len(paths) == 0 {
		return
	}
	switch target.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice:
	default:
		panic("invalid target, must be pointer,array,slice")
	}
	if target.IsNil() {
		return
	}
	pathInfo.Owner = target.Interface()

	current := paths[0]
	fv := elem(target).FieldByName(current.Name())
	if !fv.IsValid() {
		return
	}
	pathInfo.fieldValue = fv
	field := fieldValueAny(fv)

	switch p := current.(type) {
	case *pathList:
		if p.index >= fv.Len() {
			return
		}
		fv = fv.Index(p.index)
		if !fv.IsValid() {
			return
		}
		pathInfo.fieldValue = fv
		field = fieldValueAny(fv)
	case *pathListAll:
		// all index match, expand to pathLists and execute.
		for i := 0; i < fv.Len(); i++ {
			pl := &pathList{
				path: path{
					parent: current.Parent(),
					name:   current.Name(),
					ptype:  PathTypeValue,
				},
				index: i,
			}
			newPaths := make([]Path, 0, len(paths))
			newPaths = append(newPaths, pl)
			newPaths = append(newPaths, paths[1:]...)
			each(target, newPaths, pathInfo, fn)
		}
		return
	}

	// execute function, when last path element
	if len(paths) == 1 && current.Type() == PathTypeValue {
		fn(field, pathInfo)
		return
	}

	var nextTarget reflect.Value
	switch fv.Kind() {
	case reflect.Ptr:
		if fv.IsNil() {
			return
		}
		nextTarget = fv
	default:
		nextTarget = fv.Addr()
	}
	each(nextTarget, paths[1:], pathInfo, fn)
}

func fieldValueAny(fv reflect.Value) any {
	fv = reflect.Indirect(fv)
	switch {
	case !fv.IsValid():
		return nil
	case fv.CanInt():
		return fieldValueInt(fv)
	case fv.CanUint():
		return fieldValueUint(fv)
	case fv.CanFloat():
		return fieldValueFloat(fv)
	case fv.CanComplex():
		return fieldValueComplex(fv)
	case fv.Kind() == reflect.String:
		return fv.String()
	case fv.CanAddr():
		return fv.Addr()
	case fv.CanInterface():
		return fv.Interface()
	}
	return nil
}

func fieldValueInt(fv reflect.Value) any {
	n := fv.Int()
	switch fv.Kind() {
	case reflect.Int:
		return int(n)
	case reflect.Int8:
		return int8(n)
	case reflect.Int16:
		return int16(n)
	case reflect.Int32:
		return int32(n)
	case reflect.Int64:
		return n
	default:
		return n
	}
}

func fieldValueUint(fv reflect.Value) any {
	n := fv.Uint()
	switch fv.Kind() {
	case reflect.Uint:
		return uint(n)
	case reflect.Uint8:
		return uint8(n)
	case reflect.Uint16:
		return uint16(n)
	case reflect.Uint32:
		return uint32(n)
	case reflect.Uint64:
		return n
	default:
		return n
	}
}

func fieldValueFloat(fv reflect.Value) any {
	f := fv.Float()
	switch fv.Kind() {
	case reflect.Float32:
		return float32(f)
	case reflect.Float64:
		return f
	default:
		return f
	}
}

func fieldValueComplex(fv reflect.Value) any {
	c := fv.Complex()
	switch fv.Kind() {
	case reflect.Complex64:
		return complex64(c)
	case reflect.Complex128:
		return c
	default:
		return c
	}
}

func elem(v reflect.Value) reflect.Value {
	v = reflect.Indirect(v)
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}
