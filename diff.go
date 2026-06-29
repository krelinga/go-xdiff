package diff

func Diff(left, right any, differ Differ, reporter Reporter) (same bool, err error) {
	s := &State{
		Reporter: reporter,
	}
	s.Push(RootKey{}, left, right)
	return differ.Diff(s, left, right)
}
