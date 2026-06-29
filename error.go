package diff

import "fmt"

// Error is used for the error return of Differ.Diff().
//
// It wraps the original error and also contains the path at which the error happened.
type Error struct {
	Path Path
	Err  error
}

func WrapError(p Path, e error) error {
	if e == nil {
		return nil
	}
	return Error{
		Path: p,
		Err:  e,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Err.Error())
}

func (e Error) Unwrap() error {
	return e.Err
}
