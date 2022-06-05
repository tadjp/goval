package goval

import "reflect"

// Get get a value from a structure using a path.
func Get(target any, path string) any {
	p := parse(path)
	return get(target, p)
}

func get(target any, path Path) any {
	if len(path) == 0 {
		return target
	}
	v := reflect.ValueOf(target)
	p := path[0]

	fv := v.FieldByName(p)
	dist := fv.Interface()
	return get(dist, path[1:])
}
