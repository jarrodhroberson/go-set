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

I am also trying to implement some ideas from Alan Kay.

"OOP to me means only messaging, local retention and protection and
hiding of state-process, and extreme late-binding of all things."

That is why there is a Set interface with only one exported function, the one that lets you retrieve the `Set` as a `slice`.
That is why the `set` struct that tracks state nothing is exported.
Late binding of all things is really hard in a staticly typed language, so we have to do with Generics.
*/
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

// ToSlice returns a slice of values in non-deterministic order
func (s *set[T]) ToSlice() []T {
	var values []T
	for idx := range s.idValue {
		values = append(values, s.idValue[idx])
	}
	return values
}
