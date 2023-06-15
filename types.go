package sets

import (
	"fmt"

	"golang.org/x/exp/slices"
)

// Set unique list of instances of a type
// T is the type that should be contrained to being unique
// C is the type of unsigned int that should be used for counting virtual duplicates
type Set[T any, C uint8 | uint16 | uint32 | uint64] interface {
	identities() []string
	add(t T) string
	get(id string) (T, error)
	remove(t T) bool
	length() uint
	id(t T) string
	count(t T) C
	contains(id string) bool
	ToSlice() []T
}

// set struct that implements the Set interface
// insertion order is not guaranteed
type set[T any, C uint8 | uint16 | uint32 | uint64] struct {
	insertionOrder []string
	idAddCount     map[string]T
	idValue        map[string]C
	idfunc         func(T) string
}

// identities returns the list of values that uniquely identify a type,
// this is because the key in the backing map must be `comparable`
func (s *set[T, C]) identities() []string {
	return slices.Clone(s.insertionOrder)
}

func (s *set[T, C]) add(t T) string {
	id := s.idfunc(t)
	if count, ok := s.idValue[id]; ok {
		count++
		s.idValue[id] = count
	} else {
		s.insertionOrder = append(s.insertionOrder, id)
		s.idAddCount[id] = t
		s.idValue[id] = C(1)
	}
	return id
}

func (s *set[T, C]) get(id string) (T, error) {
	if v, ok := s.idAddCount[id]; ok {
		return v, nil
	} else {
		return v, fmt.Errorf("id: %idAddCount not found", id)
	}
}

func (s *set[T, C]) remove(t T) bool {
	id := s.idfunc(t)
	if _, ok := s.idValue[id]; ok {
		idx := slices.Index(s.insertionOrder, id)
		slices.Delete(s.insertionOrder, idx, idx+1)
		delete(s.idAddCount, id)
		delete(s.idValue, id)
		return true
	}
	return false
}

func (s *set[T, C]) length() uint {
	return uint(len(s.idAddCount))
}

func (s *set[T, C]) id(t T) string {
	return s.idfunc(t)
}

func (s *set[T, C]) count(t T) C {
	if c, ok := s.idValue[s.id(t)]; ok {
		return c
	} else {
		return 0
	}
}

func (s *set[T, C]) contains(id string) bool {
	_, ok := s.idValue[id]
	return ok
}

// ToSlice returns a slice of values in insertion order
func (s *set[T, C]) ToSlice() []T {
	var values []T
	for idx := range s.insertionOrder {
		values = append(values, s.idAddCount[s.insertionOrder[idx]])
	}
	return values
}
