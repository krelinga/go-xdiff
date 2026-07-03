package diff

import "reflect"

// Serial is a Task that executes the comparison serially and reports the results to a Reporter.
type Serial struct {
	// If Reporter is nil, no reporting will be done.
	Reporter Reporter2

	err error
	same bool
}

func (s *Serial) ok() bool {
	return s.err == nil
}

func (s *Serial) diff(leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	if !s.ok() {
		return
	}
	if !left.Ok() || !right.Ok() {
		panic("diff called with non-OK entries")
	}

	plan := differ.Diff(left.Must(), right.Must())
	plan.runPlan(s, leftPath, left, rightPath, right, differ)
}

func (s *Serial) leaf(leaf Leaf, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
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
			s.Reporter.Report(differ2Name(differ), leftPath, left, rightPath, right)
		}
	}
}

// TODO: delete old differName function and rename this to differName.
func differ2Name(differ Differ2) string {
	value := reflect.ValueOf(differ)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	return value.Type().String()
}

func (s *Serial) branch(branch Branch, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	if !s.ok() {
		return
	}
	s.err = branch(func(leftChildKey Key, leftChild Entry, rightChildKey Key, rightChild Entry, childDiffer Differ2) {
		if !s.ok() {
			return
		}
		leftChildPath := leftPath.Push(leftChildKey)
		rightChildPath := rightPath.Push(rightChildKey)
		if leftChild.Ok() && rightChild.Ok() {
			// Both are OK so we need to recurse into the differ.
			s.diff(leftChildPath, leftChild, rightChildPath, rightChild, childDiffer)
		} else if leftChild.Ok() || rightChild.Ok() {
			// Only one of the entries is OK so we report it as a difference.
			s.same = false
			if s.Reporter != nil {
				s.Reporter.Report(differ2Name(differ), leftChildPath, leftChild, rightChildPath, rightChild)
			}
		} else {
			// Neither entry is OK, which should not happen.
			// TODO: consider reporting this as an error instead?
			panic("both left and right entries are not OK")
		}
	})
}

func (s *Serial) compose(compose Compose, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	if !s.ok() {
		return
	}
	s.err = compose(func(childDiffer Differ2) {
		if !s.ok() {
			return
		}
		leftChildPath := leftPath.Push(composeKey{})
		rightChildPath := rightPath.Push(composeKey{})
		s.diff(leftChildPath, left, rightChildPath, right, childDiffer)
	})
}

func (s *Serial) delegate(leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	if !s.ok() {
		return
	}
	s.diff(leftPath, left, rightPath, right, differ)
}

func (s *Serial) result() (same bool, err error) {
	return s.same, s.err
}