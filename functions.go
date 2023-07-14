package sets

import (
	"github.com/jarrodhroberson/destruct"
	"golang.org/x/exp/slices"
)

// New creates a new Set using the identity.HashStructIdentity function for determining uniqueness of non-comparable T any
func New[T any](s ...T) Set[T] {
	ns := &set[T]{
		idValue: make(map[string]T, len(s)),
		idfunc:  destruct.HashIdentity[T],
	}
	Add[T](ns, s...)
	return ns
}

// Add this mutates the internal state of the Set s
func Add[T any](s Set[T], ts ...T) {
	for _, t := range ts {
		s.add(t)
	}
}

// Remove this mutates the internal state of the Set
func Remove[T any](s Set[T], ts ...T) {
	for _, t := range ts {
		s.remove(t)
	}
}

// ContainsExactly tests to see if the Sets contain exactly the same things in any order.
// This can be used as a basic "equals" test as well.
func ContainsExactly[T any](s1 Set[T], s2 Set[T]) bool {
	s1ids := slices.Clone(s1.identities())
	s2ids := slices.Clone(s2.identities())
	slices.Sort(s1ids)
	slices.Sort(s2ids)
	return slices.Equal(s1ids, s2ids)
}

// Union returns a new Set with all the items from all the Sets
func Union[T any](s1 ...Set[T]) Set[T] {
	ns := New[T]()
	for _, s := range s1 {
		Add(ns, s.ToSlice()...)
	}
	return ns
}

// Intersection returns a new Set with only the items that exist in both Sets
func Intersection[T any](s1 Set[T], s2 Set[T]) Set[T] {
	ns := New[T]()
	for _, id := range s1.identities() {
		if v, err := s2.get(id); err == nil {
			ns.add(v)
		}
	}
	return ns
}
