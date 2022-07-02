package goval_test

import (
	"fmt"
	"github.com/tadjp/goval"
	"reflect"
	"strings"
	"testing"
)

func TestSetFunc(t *testing.T) {
	type args struct {
		target any
		path   string
		newVal any
	}
	type test struct {
		name string
		args args
		want any
	}
	defaultTest := func(fn func(tt test) test) test {
		tt := test{}
		return fn(tt)
	}

	tests := []test{
		defaultTest(func(tt test) test {
			tt.name = "update field value"
			type S struct {
				V string
			}
			tt.args = args{
				target: &S{
					V: "foo",
				},
				path:   "V",
				newVal: "bar",
			}

			tt.want = &S{
				V: "bar",
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "update slice field values"
			type S struct {
				V []string
			}
			tt.args = args{
				target: &S{
					V: []string{
						"foo",
						"a",
						"b",
					},
				},
				path:   "V[0]",
				newVal: "bar",
			}

			tt.want = &S{
				V: []string{
					"bar",
					"a",
					"b",
				},
			}
			return tt
		}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, _ := goval.Parse(tt.args.path)
			goval.SetFunc[any](tt.args.target, path, func(v any, pathInfo goval.PathInfo) any {
				return tt.args.newVal
			})
			if !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("SetFunc() = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	type args struct {
		target any
		path   string
		newVal any
	}
	type test struct {
		name string
		args args
		want any
	}
	defaultTest := func(fn func(tt test) test) test {
		tt := test{}
		return fn(tt)
	}

	tests := []test{
		defaultTest(func(tt test) test {
			tt.name = "update field value"
			type S struct {
				V string
			}
			tt.args = args{
				target: &S{
					V: "foo",
				},
				path:   "V",
				newVal: "bar",
			}
			tt.want = &S{
				V: "bar",
			}
			return tt
		}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, _ := goval.Parse(tt.args.path)
			goval.Set(tt.args.target, path, tt.args.newVal)
			if !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("SetFunc() = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func ExampleSet() {
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
	goval.Set[string](&team, path, "TEAM-B")
	fmt.Println(team.Name) // TEAM-B
	// Output:
	// TEAM-B
}

func ExampleSetFunc() {
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
	path, _ := goval.Parse("Members[*].Name")
	goval.SetFunc[string](&team, path, func(v string, pathInfo goval.PathInfo) string {
		return strings.ToLower(v)
	})
	fmt.Println(team.Members[0].Name) // alice
	fmt.Println(team.Members[1].Name) // bob
	// Output:
	// alice
	// bob
}
