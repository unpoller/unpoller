package testutil

type Set[K comparable] struct {
	entities map[K]any
}

func NewSetFromMap[K comparable, V any](m map[K]V) *Set[K] {
	entities := make(map[K]any, 0)

	for k := range m {
		entities[k] = true
	}

	return &Set[K]{
		entities: entities,
	}
}

func NewSetFromSlice[K comparable](s []K) *Set[K] {
	entities := make(map[K]any, 0)

	for _, k := range s {
		entities[k] = true
	}

	return &Set[K]{
		entities: entities,
	}
}

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

func (s *Set[K]) Slice() []K {
	ret := make([]K, 0)
	for k := range s.entities {
		ret = append(ret, k)
	}

	return ret
}
