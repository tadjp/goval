package goval

import "reflect"

// Get get a value from a structure using a path.
func Get(target any, path string) any {
	var ret any
	Each(target, path, func(v any, _ any) {
		ret = v
	})
	return ret
}

// EachFunc a callback func given to each function.
//
// v: Value of the field specified in the path.
// owner: Structure owning the value.
type EachFunc func(v any, owner any)

// Each executes the given function once for each field specified in the path.
func Each(target any, path string, fn EachFunc) {
	paths := parse(path)
	each(nil, target, paths, fn)
}

func each(owner any, target any, relPath Path, fn EachFunc) {
	if len(relPath) == 0 {
		fn(target, owner)
		return
	}
	v := reflect.ValueOf(target)
	p := relPath[0]

	fv := v.FieldByName(p)
	dist := fv.Interface()
	each(target, dist, relPath[1:], fn)
}
