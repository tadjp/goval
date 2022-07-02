package goval_test

import (
	"github.com/tadjp/goval"
	"math"
	"reflect"
	"testing"
)

func TestEach(t *testing.T) {
	type args struct {
		target any
		path   string
	}
	type want struct {
		v        any
		pathInfo goval.PathInfo
	}
	type test struct {
		name    string
		args    args
		wants   []want
		wantErr bool
	}

	// create basic test data
	defaultTest := func(fn func(tt test) test) test {
		tt := test{}
		return fn(tt)
	}

	tests := []test{
		defaultTest(func(tt test) test {
			tt.name = "simple field"
			str := "foobar"
			tt.args.target = &struct {
				str string
			}{
				str: str,
			}
			tt.args.path = "str"
			tt.wants = []want{
				{
					v:        str,
					pathInfo: goval.PathInfo{Owner: tt.args.target},
				},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "simple ptr field"
			str := "foobar"
			tt.args.target = &struct {
				str *string
				num int
			}{
				str: &str,
			}
			tt.args.path = "str"
			tt.wants = []want{
				{
					v:        str,
					pathInfo: goval.PathInfo{Owner: tt.args.target},
				},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "nested struct"
			str := "foobar"
			type n struct {
				v string
			}
			type s struct {
				Nested n
			}
			target := s{
				Nested: n{
					v: str,
				},
			}
			tt.args.target = &target
			tt.args.path = "Nested.v"
			tt.wants = []want{
				{
					v:        str,
					pathInfo: goval.PathInfo{Owner: &target.Nested},
				},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "pointer struct field"
			str := "foobar"
			type n struct {
				v string
			}
			type s struct {
				Nested *n
			}
			target := s{
				Nested: &n{
					v: str,
				},
			}
			tt.args.target = &target
			tt.args.path = "Nested.v"
			tt.wants = []want{
				{
					v:        str,
					pathInfo: goval.PathInfo{Owner: target.Nested},
				},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "nil pointer struct field"
			type n struct {
				v string
			}
			type s struct {
				Nested *n
			}
			target := s{
				Nested: nil,
			}
			tt.args.target = &target
			tt.args.path = "Nested.v"
			tt.wants = nil
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "indexed slice field"
			type s struct {
				Strs []string
			}
			target := s{
				Strs: []string{
					"a",
					"b",
					"c",
				},
			}
			tt.args.target = &target
			tt.args.path = "Strs[1]"
			tt.wants = []want{
				{
					v: "b",
					pathInfo: goval.PathInfo{
						Owner: &target,
					},
				},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "all slice fields"
			type s struct {
				Strs []string
			}
			target := s{
				Strs: []string{
					"a",
					"b",
					"c",
				},
			}
			tt.args.target = &target
			tt.args.path = "Strs[*]"
			tt.wants = []want{
				{
					v: "a",
					pathInfo: goval.PathInfo{
						Owner: &target,
					},
				},
				{
					v: "b",
					pathInfo: goval.PathInfo{
						Owner: &target,
					},
				},
				{
					v: "c",
					pathInfo: goval.PathInfo{
						Owner: &target,
					},
				},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.name = "complex slice fields"
			type t struct {
				Vars []string
			}
			type s struct {
				Ts []t
			}
			target := s{
				Ts: []t{
					{
						Vars: []string{
							"a", "b", "c",
						},
					},
					{
						Vars: []string{
							"A", "B", "C",
						},
					},
				},
			}
			tt.args.target = &target
			tt.args.path = "Ts[*].Vars[2]"
			tt.wants = []want{
				{
					v: "c",
					pathInfo: goval.PathInfo{
						Owner: &target.Ts[0],
					},
				},
				{
					v: "C",
					pathInfo: goval.PathInfo{
						Owner: &target.Ts[1],
					},
				},
			}
			return tt
		}),
	}

	s := struct {
		Int        int
		UInt       uint
		Int8       int8
		UInt8      uint8
		Int16      int16
		UInt16     uint16
		Int32      int32
		UInt32     uint32
		Int64      int64
		UInt64     uint64
		Float32    float32
		Float64    float64
		Complex64  complex64
		Complex128 complex128
	}{
		Int:        math.MinInt,
		UInt:       math.MaxUint,
		Int8:       math.MinInt8,
		UInt8:      math.MaxUint8,
		Int16:      math.MinInt16,
		UInt16:     math.MaxUint16,
		Int32:      math.MinInt32,
		UInt32:     math.MaxUint32,
		Int64:      math.MinInt64,
		UInt64:     math.MaxUint64,
		Float32:    math.MaxFloat32,
		Float64:    math.MaxFloat64,
		Complex64:  complex(math.MaxFloat32, math.MaxFloat32),
		Complex128: complex(math.MaxFloat64, math.MaxFloat64),
	}
	tests = append(tests,
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Int"
			tt.wants = []want{
				{v: s.Int, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "UInt"
			tt.wants = []want{
				{v: s.UInt, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Int8"
			tt.wants = []want{
				{v: s.Int8, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "UInt8"
			tt.wants = []want{
				{v: s.UInt8, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Int16"
			tt.wants = []want{
				{v: s.Int16, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "UInt16"
			tt.wants = []want{
				{v: s.UInt16, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Int32"
			tt.wants = []want{
				{v: s.Int32, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "UInt32"
			tt.wants = []want{
				{v: s.UInt32, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Int64"
			tt.wants = []want{
				{v: s.Int64, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "UInt64"
			tt.wants = []want{
				{v: s.UInt64, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Float32"
			tt.wants = []want{
				{v: s.Float32, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Float64"
			tt.wants = []want{
				{v: s.Float64, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Complex64"
			tt.wants = []want{
				{v: s.Complex64, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
		defaultTest(func(tt test) test {
			tt.args.target = &s
			tt.args.path = "Complex128"
			tt.wants = []want{
				{v: s.Complex128, pathInfo: goval.PathInfo{Owner: &s}},
			}
			return tt
		}),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type got struct {
				v        any
				pathInfo goval.PathInfo
			}
			var gots []got

			path, err := goval.Parse(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Each(%v) error = %v, wantErr %v", tt.args.path, err, tt.wantErr)
			}

			goval.Each(tt.args.target, path, func(v any, pathInfo goval.PathInfo) {
				gots = append(gots, got{
					v:        v,
					pathInfo: pathInfo,
				})
			})

			if len(gots) != len(tt.wants) {
				t.Errorf("Each(%v) num of call = %v, want %v", tt.args.path, len(gots), len(tt.wants))
				return
			}
			for i, got := range gots {
				want := tt.wants[i]
				if !reflect.DeepEqual(got.v, want.v) {
					t.Errorf("Each(%v) value[%d] = %v, want %v", tt.args.path, i, got.v, want.v)
				}

				if got.pathInfo.Owner != want.pathInfo.Owner {
					t.Errorf("Each(%#v) [%d]pathInfo.Owner of value = %v, want %v",
						tt.args.path, i, got.pathInfo.Owner, want.pathInfo.Owner)
				}
			}
		})
	}
}
