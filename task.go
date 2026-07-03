package diff

type Task interface {
	leaf(leaf Leaf, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2)
	branch(branch Branch, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2)
	compose(compose Compose, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2)
	delegate(leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2)

	result() (same bool, err error)
}