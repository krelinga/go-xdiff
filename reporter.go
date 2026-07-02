package diff

// TODO: rename to Reporter and delete the old Reporter interface
type Reporter2 interface {
	Report(differName string, leftPath Path, left Entry, rightPath Path, right Entry)
}
