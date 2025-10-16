package set

type Set[T comparable] struct {
	Data map[T]struct{}
}

func NewSet[T comparable](elements ...T) *Set[T] {
	newSet := &Set[T]{
		Data: make(map[T]struct{}),
	}

	newSet.Add(elements...)

	return newSet
}

func (s *Set[T]) Add(elements ...T) {
	for _, el := range elements {
		s.Data[el] = struct{}{}
	}
}
