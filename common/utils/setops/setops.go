package setops

func Intersect[K comparable](a, b map[K]struct{}) map[K]struct{} {
	result := make(map[K]struct{})
	for k := range a {
		if _, ok := b[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// convert map[K]V to map[K]struct{}
func Mtos[K comparable, V any](m map[K]V) map[K]struct{} {
	result := make(map[K]struct{})
	for k := range m {
		result[k] = struct{}{}
	}
	return result
}