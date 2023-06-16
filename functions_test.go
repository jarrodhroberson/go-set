package sets

import (
	"sort"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestAddStrings(t *testing.T) {
	type args[T any] struct {
		s  Set[T]
		ts []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "slice of strings",
			args: args[string]{s: New[string](strings.Split("1,2,3,4,5,1,2,3,4,5", ",")...), ts: strings.Split("1,2,3,4,5", ",")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Add(tt.args.s, tt.args.ts...)
			original := tt.args.s.ToSlice()
			sort.Strings(original)
			if !slices.Equal(original, tt.args.ts) {
				t.Errorf("%v not equal to %v", original, tt.args.ts)
			}
		})
	}
}

func TestAddInts(t *testing.T) {
	type args[T any] struct {
		s  Set[T]
		ts []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "slice of ints",
			args: args[int]{s: New[int](0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9), ts: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Add(tt.args.s, tt.args.ts...)
			original := tt.args.s.ToSlice()
			sort.Ints(original)
			if !slices.Equal[int](original, tt.args.ts) {
				t.Errorf("%v not equal to %v", original, tt.args.ts)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	type args[T any] struct {
		ss []Set[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want Set[T]
	}
	tests := []testCase[int]{
		{
			name: "sets of ints",
			args: args[int]{
				ss: []Set[int]{New(0, 1, 2, 3), New(4, 5, 6, 7), New(8, 9)},
			},
			want: New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Union(tt.args.ss...); !ContainsExactly(got, tt.want) {
				t.Errorf("From() = %v, want %v", got.ToSlice(), tt.want.ToSlice())
			}
		})
	}
}
