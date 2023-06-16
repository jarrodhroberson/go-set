package sets

import (
	"fmt"
)

// Set unique list of instances of a type
// T is the type that should be contrained to being unique
// C is the type of unsigned int that should be used for counting virtual duplicates
type Set[T any] interface {
	identities() []string
	add(t T) string
	get(id string) (T, error)
	remove(t T) bool
	length() uint
	id(t T) string
	contains(id string) bool
	ToSlice() []T
}

// set struct that implements the Set interface
// insertion order is not guaranteed
type set[T any] struct {
	idValue map[string]T
	idfunc  func(T) string
}

// identities returns the list of values that uniquely identify a type,
// this is because the key in the backing map must be `comparable`
func (s *set[T]) identities() []string {
	ids := make([]string, 0, len(s.idValue))
	for k := range s.idValue {
		ids = append(ids, k)
	}
	return ids
}

func (s *set[T]) add(t T) string {
	id := s.idfunc(t)
	s.idValue[id] = t
	return id
}

func (s *set[T]) get(id string) (T, error) {
	if v, ok := s.idValue[id]; ok {
		return v, nil
	} else {
		return v, fmt.Errorf("id: %s not found", id)
	}
}

func (s *set[T]) remove(t T) bool {
	id := s.idfunc(t)
	if _, ok := s.idValue[id]; ok {
		delete(s.idValue, id)
		return true
	}
	return false
}

func (s *set[T]) length() uint {
	return uint(len(s.idValue))
}

func (s *set[T]) id(t T) string {
	return s.idfunc(t)
}

func (s *set[T]) contains(id string) bool {
	_, ok := s.idValue[id]
	return ok
}

// ToSlice returns a slice of values in non deterministic order
func (s *set[T]) ToSlice() []T {
	var values []T
	for idx := range s.idValue {
		values = append(values, s.idValue[idx])
	}
	return values
}
