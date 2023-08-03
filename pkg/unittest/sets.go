package unittest

// Set provides a generic way to compare sets of type K. This is only used for unit testing.
type Set[K comparable] struct {
	entities map[K]any
}

// NewSetFromMap will create a Set of type K from a map[K]V. V is useless here as we are only comparing the set of keys
// in the map.
func NewSetFromMap[K comparable, V any](m map[K]V) *Set[K] {
	entities := make(map[K]any, 0)

	for k := range m {
		entities[k] = true
	}

	return &Set[K]{
		entities: entities,
	}
}

// NewSetFromSlice will create a Set of type K from a slice of keys. Duplicates will be dropped as this is a set.
func NewSetFromSlice[K comparable](s []K) *Set[K] {
	entities := make(map[K]any, 0)

	for _, k := range s {
		entities[k] = true
	}

	return &Set[K]{
		entities: entities,
	}
}

// Difference will compare two this Set against another Set of the same type K. This will return entries that
// exist in this set but not the other as `additions` and entries that exist in the other set but not this set
// as `deletions`.
func (s *Set[K]) Difference(other *Set[K]) (additions []K, deletions []K) {
	additions = make([]K, 0)

	for i := range s.entities {
		if _, ok := other.entities[i]; !ok {
			additions = append(additions, i)
		}
	}

	deletions = make([]K, 0)

	for j := range other.entities {
		if _, ok := s.entities[j]; !ok {
			deletions = append(deletions, j)
		}
	}

	return additions, deletions
}

func (s *Set[K]) Len() int {
	return len(s.entities)
}

// Slice will return the set back as a slice of type K
func (s *Set[K]) Slice() []K {
	ret := make([]K, 0)
	for k := range s.entities {
		ret = append(ret, k)
	}

	return ret
}
