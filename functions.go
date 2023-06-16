/*
Package sets is a package that represents as idiomatic `Set` implemenation as I can conceive of as of the time of the last commit.
Most `Set` implementations focus on a more C++/Java/C# false "OOP" approach. Defining an `interface` that some `struct` then has a bunch of
receiver methods that manipulate the internal state of the `Set` and ensure its contractual behavior.

Receiver methods on value types, is not the most idiomatic Go. So my approach is package level functions that apply the
`Set` logic and maintain high cohesion and loose coupling from my implementation details.

This way client code does not get tightly coupled with the interface and implementation.
This is mainly because `Set` semantics are only really enforced when building or adding to the `Set`.
Thus, once you finish creating the `Set` you can easily call `ToSlice()` get an insertion order slice of the values
and used them other code that expects `[]T` and not `Set`.

There is no `SortedSet` implementation because that mixes concerns, a `Set` concerns are simple, no duplicates.
If you need a `SortedSet` then you build your `Set` and then sort the resulting `ToSlice()` results.
*/
package sets

import (
	"golang.org/x/exp/slices"

	"github.com/jarrodhroberson/go-set/internal/identity"
)

// New creates a new Set using the identity.HashStructIdentity function for determining uniqueness of non-comparable T any
func New[T any](s ...T) Set[T] {
	ns := &set[T]{
		idValue: make(map[string]T, len(s)),
		idfunc:  identity.HashIdentity[T],
	}
	Add[T](ns, s...)
	return ns
}

func From[T any](s []T) Set[T] {
	return New[T](s...)
}

// Add this mutates the internal state of the Set
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
	return slices.Equal(s1.identities(), s2.identities())
}

// Union returns a new Set with all the items from all the Sets
func Union[T any](s1 ...Set[T]) Set[T] {
	ns := New[T]()
	for _, s := range s1 {
		Add(New[T](s.ToSlice()...))
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
