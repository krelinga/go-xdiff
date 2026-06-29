package diff

import "strings"

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

type Reporter interface {
	// Called when starting a new layer in the nested hierarchy.
	Push(key Key, left, right any)
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
	Path Path
}

func (s *State) Push(key Key, left, right any) {
	s.Path = append(s.Path, key)
	if s.Reporter != nil {
		s.Reporter.Push(key, left, right)
	}
}

func (s *State) Pop() {
	s.Path = s.Path[0:len(s.Path)-1]
	if s.Reporter != nil {
		s.Reporter.Pop()
	}
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