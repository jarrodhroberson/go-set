package sets

import (
	"fmt"
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
	s      map[string]T
	c      map[string]C
	idfunc func(T) string
}

// identities returns the list of values that uniquely identify a type,
// this is because the key in the backing map must be `comparable`
func (s set[T, C]) identities() []string {
	var ids = make([]string, 0, len(s.s))
	for k := range s.s {
		ids = append(ids, k)
	}
	return ids
}

func (s set[T, C]) add(t T) string {
	id := s.idfunc(t)
	if count, ok := s.c[id]; ok {
		count++
		s.c[id] = count
	} else {
		s.s[id] = t
		s.c[id] = C(1)
	}
	return id
}

func (s set[T, C]) get(id string) (T, error) {
	if v, ok := s.s[id]; ok {
		return v, nil
	} else {
		return v, fmt.Errorf("id: %s not found")
	}
}

func (s set[T, C]) remove(t T) bool {
	id := s.idfunc(t)
	if _, ok := s.c[id]; ok {
		delete(s.s, id)
		delete(s.c, id)
		return true
	}
	return false
}

func (s set[T, C]) length() uint {
	return uint(len(s.s))
}

func (s set[T, C]) id(t T) string {
	return s.idfunc(t)
}

func (s set[T, C]) count(t T) C {
	if c, ok := s.c[s.id(t)]; ok {
		return c
	} else {
		return 0
	}
}

func (s set[T, C]) contains(id string) bool {
	_, ok := s.c[id]
	return ok
}

func (s set[T, C]) ToSlice() []T {
	var values []T
	for k := range s.s {
		values = append(values, s.s[k])
	}
	return values
}
