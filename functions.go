/*
Package sets is a package that represents as idiomatic `Set` implemenation as I can conceive of as of the time of the last commit.
Most `Set` implementations focus on a more C++/Java/C# false "OOP" approach. Defining an `interface` that some `struct` then has a bunch of
receiver methods that manipulate the internal state of the `Set` and ensure its contractual behavior.

Receiver methods on value types, is not the most idiomatic Go. So my approach is package level functions that apply the
`Set` logic and maintain high cohesion and loose coupling from my implementation details.

This way client code does not get tightly coupled with the interface and implementation.
This is mainly because `Set` semantics are only really inforced when building or adding to the `Set`.
Thus, once you finish creating the `Set` you can easily call `ToSlice()` get an insertion order slice of the values
and used them other code that expects `[]T` and not `Set`.

There is no `SortedSet` implementation because that mixes concerns, a `Set` concerns are simple, no duplicates.
If you need a `SortedSet` then you build your `Set` and then sort the resulting `ToSlice()` results.

I only added tracking copies of duplicates added because the value portion of that map was going to be empty any way.
Might as well put it to some use and keep track of how many duplicates were added. This would allow one to use it as a
sparse slice and generate the duplicates back again if needed.
*/
package sets

import (
	"golang.org/x/exp/slices"

	"github.com/jarrodhroberson/go-set/internal/identity"
)

// New creates a new Set using the identity.HashStructIdentity function for determining uniqueness of non-comparable T any
func New[T any, C uint8 | uint16 | uint32 | uint64](s ...T) Set[T, C] {
	ns := &set[T, C]{
		insertionOrder: make([]string, 9, len(s)),
		idAddCount:     make(map[string]T, len(s)),
		idValue:        make(map[string]C, len(s)),
		idfunc:         identity.HashIdentity[T],
	}
	Add[T, C](ns, s...)
	return ns
}

// Add this mutates the internal state of the Set
func Add[T any, C uint8 | uint16 | uint32 | uint64](s Set[T, C], ts ...T) {
	for _, t := range ts {
		s.add(t)
	}
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
