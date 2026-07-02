package diff

type Task interface {
	leaf(Leaf)
	branch(Branch)
	compose(Compose)
	delegate(Differ)

	result() (same bool, err error)
}