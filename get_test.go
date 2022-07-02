package goval_test

import (
	"fmt"
	"reflect"

	"github.com/tadjp/goval"

	"testing"
)

func TestGetAll(t *testing.T) {
	type address struct {
		CountryCode int
		PostalCode  string
	}

	type person struct {
		Name       string
		Age        int
		Address    address
		PtrAddress *address
	}

	type args struct {
		src  any
		path string
	}
	type test struct {
		name string
		args args
		want []any
	}

	// create basic test data
	defaultTest := func(fn func(tt test) test) test {
		tt := test{
			args: args{
				src: &person{
					Name: "Alice",
					Age:  25,
					Address: address{
						CountryCode: 81,
						PostalCode:  "1000001",
					},
					PtrAddress: &address{
						CountryCode: 1,
						PostalCode:  "0000002",
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
			tt.want = []any{"Alice"}
			return tt
		}),

		defaultTest(func(tt test) test {
			tt.name = "get a nested field value"
			tt.args.path = "Address.CountryCode"
			tt.want = []any{81}
			return tt
		}),

		defaultTest(func(tt test) test {
			tt.name = "get a pointer field value"
			tt.args.path = "PtrAddress.PostalCode"
			tt.want = []any{"0000002"}
			return tt
		}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, _ := goval.Parse(tt.args.path)
			got := goval.GetAll[any](tt.args.src, path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll(%v) = %v, want %v", tt.args.path, got, tt.want)
			}
		})
	}
}

func ExampleGetAll() {
	type Member struct {
		Name string
	}

	type Team struct {
		Name    string
		Members []*Member
	}

	team := Team{
		Name: "TEAM-A",
		Members: []*Member{
			{
				Name: "Alice",
			},
			{
				Name: "Bob",
			},
		},
	}
	path, _ := goval.Parse("Name")
	fmt.Println(goval.GetAll[string](&team, path)) // [TEAM-A]
	path, _ = goval.Parse("Members[0].Name")
	fmt.Println(goval.GetAll[string](&team, path)) // [Alice]
	path, _ = goval.Parse("Members[*].Name")
	fmt.Println(goval.GetAll[string](&team, path)) // [Alice Bob]
	// Output:
	// [TEAM-A]
	// [Alice]
	// [Alice Bob]
}
