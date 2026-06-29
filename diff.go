package diff

func Diff(left, right any, differ Differ, reporter Reporter) (same bool, err error) {
	s := &State{
		Reporter: reporter,
	}
	return s.DiffChild(RootKey{}, left, right, differ)
}
