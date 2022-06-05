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

func TestEach(t *testing.T) {
	type args struct {
		target any
		path   string
	}
	type test struct {
		name      string
		args      args
		wantVal   any
		wantOwner any
	}
	// create basic test data
	defaultTest := func(fn func(tt test) test) test {
		tt := test{
			args: args{
				target: person{
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
			tt.name = "each a field value"
			tt.args.path = "Name"

			tt.wantVal = "Alice"
			tt.wantOwner = tt.args.target
			return tt
		}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotVal any
			var gotOwner any
			godot.Each(tt.args.target, tt.args.path, func(v any, owner any) {
				gotVal = v
				gotOwner = owner
			})

			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Each(%v) path's value = %v, want %v", tt.args.path, gotVal, tt.wantVal)
			}
			if !reflect.DeepEqual(gotOwner, tt.wantOwner) {
				t.Errorf("Each(%v) owner of value = %v, want %v", tt.args.path, gotOwner, tt.wantOwner)
			}
		})
	}
}
