package sets

import (
	"github.com/jarrodhroberson/go-set/identity"
	"golang.org/x/exp/slices"
)

// New creates a new Set using the identity.HashStructIdentity function for determining uniqueness of non-comparable T any
func New[T any, C uint8 | uint16 | uint32 | uint64](s ...T) Set[T, C] {
	ns := set[T, C]{
		idfunc: identity.HashStructIdentity[T],
	}
	_ = Add[T, C](ns, s...)
	return ns
}

// Add this mutates the internal state of the Set
func Add[T any, C uint8 | uint16 | uint32 | uint64](s Set[T, C], ts ...T) []string {
	ids := make([]string, 0, len(ts))
	for _, t := range ts {
		ids = append(ids, s.add(t))
	}
	return ids
}

// Remove this mutates the internal state of the Set
func Remove[T any, C uint8 | uint16 | uint32 | uint64](s Set[T, C], ts ...T) {
	for _, t := range ts {
		s.remove(t)
	}
}

// ContainsExactly tests to see if the Sets contain exactly the same things in any order.
// This can be used as a basic "equals" test as well.
func ContainsExactly[T any, C uint8 | uint16 | uint32 | uint64](s1 Set[T, C], s2 Set[T, C]) bool {
	s1ids := slices.Clone(s1.identities())
	s2ids := slices.Clone(s2.identities())
	slices.Sort(s1ids)
	slices.Sort(s2ids)
	return slices.Equal(s1.identities(), s2.identities())
}

// Union returns a new Set with all the items from all the Sets
func Union[T any, C uint8 | uint16 | uint32 | uint64](s1 ...Set[T, C]) Set[T, C] {
	ns := New[T, C]()
	for _, s := range s1 {
		Add(New[T, C](s.ToSlice()...))
	}
	return ns
}

// Intersection returns a new Set with only the items that exist in both Sets
func Intersection[T any, C uint8 | uint16 | uint32 | uint64](s1 Set[T, C], s2 Set[T, C]) Set[T, C] {
	ns := New[T, C]()
	for _, id := range s1.identities() {
		if v, err := s2.get(id); err == nil {
			ns.add(v)
		}
	}
	return ns
}
