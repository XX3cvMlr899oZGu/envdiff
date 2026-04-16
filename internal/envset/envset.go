// Package envset provides utilities for set operations on env maps.
package envset

// Intersection returns keys present in all provided maps.
func Intersection(maps ...map[string]string) map[string]string {
	if len(maps) == 0 {
		return map[string]string{}
	}
	result := make(map[string]string)
	for k, v := range maps[0] {
		result[k] = v
	}
	for _, m := range maps[1:] {
		for k := range result {
			if _, ok := m[k]; !ok {
				delete(result, k)
			}
		}
	}
	return result
}

// Union returns all keys from all provided maps.
// Later maps overwrite earlier ones on conflict.
func Union(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Difference returns keys present in base but not in any of the others.
func Difference(base map[string]string, others ...map[string]string) map[string]string {
	exclude := make(map[string]struct{})
	for _, m := range others {
		for k := range m {
			exclude[k] = struct{}{}
		}
	}
	result := make(map[string]string)
	for k, v := range base {
		if _, found := exclude[k]; !found {
			result[k] = v
		}
	}
	return result
}
