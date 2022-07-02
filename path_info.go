package goval

import "reflect"

type PathInfo struct {
	RequirePath Path
	Owner       any //
	fieldValue  reflect.Value
}
