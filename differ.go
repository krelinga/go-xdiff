package diff

import (
	"reflect"
	"slices"
	"strings"
)

type Key interface {
	DiffKey() string
}

type Path []Key

func (p Path) String() string {
	parts := make([]string, len(p))
	for i, key := range p {
		parts[i] = key.DiffKey()
	}
	return strings.Join(parts, " > ")
}

func PathEqual(a, b Path) bool {
	return slices.Equal(a, b)
}

type Reporter interface {
	// Called when starting a new layer in the nested hierarchy.
	Push(key Key, name string, left, right any)
	// Called when finished with a layer in the nested hierarchy.
	Pop()

	// Called to report an entry that is only present in the left side of the diff.
	LeftOnly(key Key, left any)
	// Called to report an entry that is only present in the right side of the diff.
	RightOnly(key Key, right any)

	// Called to report that a difference exists in the current layer of the hierarchy.
	Different()
}

type State struct {
	Reporter Reporter
	Path     Path
}

func (s *State) push(key Key, name string, left, right any) {
	s.Path = append(s.Path, key)
	if s.Reporter != nil {
		s.Reporter.Push(key, name, left, right)
	}
}

func (s *State) pop() {
	s.Path = s.Path[0 : len(s.Path)-1]
	if s.Reporter != nil {
		s.Reporter.Pop()
	}
}

func differName(differ Differ) string {
	value := reflect.ValueOf(differ)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value.Type().String()
}

func (s *State) DiffChild(key Key, left, right any, differ Differ) (same bool, err error) {
	s.push(key, differName(differ), left, right)
	defer s.pop()
	return differ.Diff(s, left, right)
}

func (s *State) LeftOnly(key Key, left any) {
	if s.Reporter != nil {
		s.Reporter.LeftOnly(key, left)
	}
}

func (s *State) RightOnly(key Key, right any) {
	if s.Reporter != nil {
		s.Reporter.RightOnly(key, right)
	}
}

func (s *State) Different() {
	if s.Reporter != nil {
		s.Reporter.Different()
	}
}

// Differs encapsulate the logic to diff various kinds of data.
type Differ interface {
	// state will never be nil.
	Diff(state *State, left, right any) (same bool, err error)
}
