package identity

import (
	"testing"
	"time"
)

type person struct {
	FirstName    string
	LastName     string
	EmailAddress string
	BirthDate    time.Time
	CreatedDate  time.Time `identity:"-"`
}

func mustLocation(ianaName string) *time.Location {
	l, err := time.LoadLocation(ianaName)
	if err != nil {
		panic(err)
	}
	return l
}

var EST = mustLocation("America/New_York")

var me = person{
	FirstName:    "Jarrod",
	LastName:     "Roberson",
	EmailAddress: "jarrod@vertigrated.com",
	BirthDate:    time.Date(1967, time.November, 27, 0, 0, 0, 0, EST),
	CreatedDate:  time.Now().UTC(),
}

func TestHashIdentity(t *testing.T) {
	type args[T any] struct {
		t T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want string
	}
	tests := []testCase[person]{
		{
			name: "person exclude CreatedDate",
			args: args[person]{t: me},
			want: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashIdentity(tt.args.t); got != tt.want {
				t.Errorf("HashIdentity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashIdentityString(t *testing.T) {
	type args[T any] struct {
		t T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want string
	}
	tests := []testCase[string]{
		{
			name: "1",
			args: args[string]{"1"},
			want: "?",
		},
		{
			name: "2",
			args: args[string]{"1"},
			want: "?",
		},
		{
			name: "3",
			args: args[string]{"1"},
			want: "?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashIdentity(tt.args.t); got != tt.want {
				t.Errorf("HashIdentity() = %v, want %v", got, tt.want)
			}
		})
	}
}
