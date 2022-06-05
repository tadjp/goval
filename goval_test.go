package goval_test

import (
	"reflect"

	godot "github.com/tadjp/goval"

	"testing"
)

type address struct {
	CountryCode int
	PostalCode  string
}

type person struct {
	Name    string
	Age     int
	Address address
}

func TestGet(t *testing.T) {
	type args struct {
		src  any
		path string
	}
	type test struct {
		name string
		args args
		want any
	}

	// create basic test data
	defaultTest := func(fn func(tt test) test) test {
		tt := test{
			args: args{
				src: person{
					Name: "Alice",
					Age:  25,
					Address: address{
						CountryCode: 81,
						PostalCode:  "1000001",
					},
				},
			},
		}
		return fn(tt)
	}

	tests := []test{
		defaultTest(func(tt test) test {
			tt.name = "get a field value"
			tt.args.path = "Name"
			tt.want = "Alice"
			return tt
		}),

		defaultTest(func(tt test) test {
			tt.name = "get a nested field value"
			tt.args.path = "Address.CountryCode"
			tt.want = 81
			return tt
		}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := godot.Get(tt.args.src, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get(%v) = %v, want %v", tt.args.path, got, tt.want)
			}
		})
	}
}
