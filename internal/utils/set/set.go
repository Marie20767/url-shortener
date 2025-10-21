package set

type Set[T comparable] struct {
	Data map[T]struct{}
}

func New[T comparable](elements ...T) *Set[T] {
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

func (s *Set[T]) ToSlice() []T {
	var slice []T

	for key := range s.Data {
		slice = append(slice, key)
	}

	return slice
}
