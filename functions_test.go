package sets

import (
	"fmt"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	type args[T any, C interface {
		uint8 | uint16 | uint32 | uint64
	}] struct {
		s  Set[T, C]
		ts []T
	}
	type testCase[T any, C interface {
		uint8 | uint16 | uint32 | uint64
	}] struct {
		name string
		args args[T, C]
	}
	tests := []testCase[string, uint64]{
		{
			name: "1,2,3 + 4,5",
			args: args[string, uint64]{s: New[string, uint64](strings.Split("1,2,3", ",")...), ts: strings.Split("4,5", ",")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Add(tt.args.s, tt.args.ts...)
			fmt.Println(strings.Join(tt.args.s.ToSlice(), ","))
		})
	}
}
