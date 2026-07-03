package diff

import "reflect"

type serialLevel struct {
	left, right Entry
	leftKey, rightKey Key
	differ Differ2
}

// Serial is a Task that executes the comparison serially and reports the results to a Reporter.
type Serial struct {
	// If Reporter is nil, no reporting will be done.
	Reporter Reporter2

	levels []*serialLevel

	err error
	same bool
}

func (s *Serial) ok() bool {
	return s.err == nil
}

func (s *Serial) push(leftKey Key, left Entry, rightKey Key, right Entry, differ Differ2) {
	s.levels = append(s.levels, &serialLevel{
		left: left,
		right: right,
		leftKey: leftKey,
		rightKey: rightKey,
		differ: differ,
	})
}

func (s *Serial) pop() {
	if len(s.levels) == 0 {
		panic("pop called on empty stack")
	}
	s.levels = s.levels[:len(s.levels)-1]
}

func (s *Serial) top() *serialLevel {
	if len(s.levels) == 0 {
		panic("top called on empty stack")
	}
	return s.levels[len(s.levels)-1]
}

func (s *Serial) differName() string {
	differ := s.top().differ
	if differ == nil {
		panic("differ is nil")
	}
	return reflect.TypeOf(differ).String()
}

func (s *Serial) leftPath() Path {
	path := Path(make([]Key, 0, len(s.levels)))
	for _, level := range s.levels {
		if level.leftKey == nil {
			continue
		}
		path = append(path, level.leftKey)
	}
	return path
}

func (s *Serial) rightPath() Path {
	path := Path(make([]Key, 0, len(s.levels)))
	for _, level := range s.levels {
		if level.rightKey == nil {
			continue
		}
		path = append(path, level.rightKey)
	}
	return path
}

func (s *Serial) leftEntry() Entry {
	return s.top().left
}

func (s *Serial) rightEntry() Entry {
	return s.top().right
}

func (s *Serial) run() {
	if !s.ok() {
		return
	}
	top := s.top()
	top.differ.Diff(top.left, top.right).runPlan(s)
}

func (s *Serial) leaf(leaf Leaf) {
	if !s.ok() {
		return
	}
	var same bool
	same, s.err = leaf()
	if !s.ok() {
		return
	}
	if !same {
		s.same = false
		if s.Reporter != nil {
			s.Reporter.Report(s.differName(), s.leftPath(), s.leftEntry(), s.rightPath(), s.rightEntry())
		}
	}
}

func (s *Serial) branch(branch Branch) {
	if !s.ok() {
		return
	}
	s.err = branch(func(leftKey Key, left Entry, rightKey Key, right Entry, differ Differ2) {
		if !s.ok() {
			return
		}
		if left.Ok() && right.Ok() {
			// Both are OK so we need to recurse into the differ.
			s.push(leftKey, left, rightKey, right, differ)
			defer s.pop()
			s.run()
		} else if left.Ok() || right.Ok() {
			// Only one of the entries is OK so we report it as a difference.
			s.same = false
			if s.Reporter != nil {
				s.Reporter.Report(s.differName(), s.leftPath(), s.leftEntry(), s.rightPath(), s.rightEntry())
			}
		} else {
			// Neither entry is OK, which should not happen.
			// TODO: consider reporting this as an error instead?
			panic("both left and right entries are not OK")
		}
	})
}

func (s *Serial) compose(compose Compose) {
	if !s.ok() {
		return
	}
	s.err = compose(func(differ Differ2) {
		if !s.ok() {
			return
		}
		s.push(nil, s.leftEntry(), nil, s.rightEntry(), differ)
		defer s.pop()
		s.run()
	})
}

func (s *Serial) delegate(differ Differ2) {
	if !s.ok() {
		return
	}
	oldDiffer := s.top().differ
	s.top().differ = differ
	defer func() {
		s.top().differ = oldDiffer
	}()
	s.run()
}

func (s *Serial) result() (same bool, err error) {
	return s.same, s.err
}