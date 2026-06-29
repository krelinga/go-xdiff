package diff

// Ignore is a Differ that always reports no differences and no errors.
type Ignore struct{}

func (_ Ignore) Diff(_ *State, _, _ any) (bool, error) {
	return false, nil
}