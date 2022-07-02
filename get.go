package goval

// GetAll get field values
func GetAll[T any](target any, path Path) []T {
	s := make([]T, 0)
	Each(target, path, func(v any, pathInfo PathInfo) {
		r, ok := v.(T)
		if !ok {
			panic("invalid type assign")
		}
		s = append(s, r)
	})
	return s
}
