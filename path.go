package goval

import (
	"errors"
	"strconv"
	"strings"
)
import "regexp"

// Path
//
// example.
// p = NewPath("/person/age")
// p.Name() // "age"
// p.Parent().Name() // "person"
// p.Type() // PathTypeValue
type Path interface {
	Parent() Path
	Name() string
	Split() []Path
	Type() PathType
}

type PathType int

const (
	_ PathType = iota
	PathTypeValue
	PathTypeCollection
)

// Parse parsing path.
func Parse(pathStr string) (_ Path, err error) {
	tokens := strings.Split(pathStr, ".")
	var p Path
	for _, token := range tokens {
		p, err = parseToken(p, token)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

var regList *regexp.Regexp
var regValue *regexp.Regexp

func init() {
	var err error
	if regList, err = regexp.Compile(`(\w+)\[(\d+|\*)?]`); err != nil {
		panic(err)
	}
	if regValue, err = regexp.Compile(`\w+`); err != nil {
		panic(err)
	}
}

func parseToken(parent Path, str string) (Path, error) {
	switch {
	case regList.MatchString(str): // match list path. e.g. foo.bar[*]
		group := regList.FindStringSubmatch(str)
		name := group[1]
		idxStr := group[2]
		switch idxStr {
		case "", "*":
			return &pathListAll{
				path: path{
					parent: parent,
					name:   name,
					ptype:  PathTypeCollection,
				},
				all: true,
			}, nil
		default:
			i, err := strconv.ParseInt(idxStr, 10, 32)
			if err != nil {
				return nil, err
			}
			return &pathList{
				path: path{
					parent: parent,
					name:   name,
					ptype:  PathTypeValue,
				},
				index: int(i),
			}, nil
		}
	case regValue.MatchString(str): // match value path. e.g. foo.Name
		return &path{
			parent: parent,
			ptype:  PathTypeValue,
			name:   str,
		}, nil
	}

	return nil, errors.New("invalid defined path")
}

type path struct {
	parent Path
	name   string
	ptype  PathType
}

func (p *path) Parent() Path {
	return p.parent
}

func (p *path) Name() string {
	return p.name
}

func (p *path) Split() []Path {
	return splitPath(p)
}

func (p *path) Type() PathType {
	return p.ptype
}

type pathList struct {
	path
	index int
}

func (p *pathList) Split() []Path {
	return splitPath(p)
}

type pathListAll struct {
	path
	all bool
}

func (p *pathListAll) Split() []Path {
	return splitPath(p)
}

func splitPath(p Path) []Path {
	if p.Parent() == nil {
		return []Path{p}
	}
	return append(splitPath(p.Parent()), p)
}
