package diff

type Task interface {
	leaf(Leaf)
	branch(Branch)
	compose(Compose)
	delegate(Differ2)

	result() (same bool, err error)
}