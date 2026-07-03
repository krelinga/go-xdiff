package diff

type Plan interface {
	runPlan(task Task, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2)
}

// Leaf is a Plan that compares a single leaf field across left and right.
type Leaf func() (same bool, err error)

func (l Leaf) runPlan(task Task, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	task.leaf(l, leftPath, left, rightPath, right, differ)
}

// YieldEntry is a function that is used to yield a single matched set of entries from a branch.
type YieldEntry = func(leftKey Key, left Entry, rightKey Key, right Entry, differ Differ2)

type Branch func(YieldEntry) error

func (b Branch) runPlan(task Task, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	task.branch(b, leftPath, left, rightPath, right, differ)
}

type YieldDiffer = func(Differ2)

// Compose is a Plan that allows a Differ to be implemented in terms of one or more other Differ implementations.
type Compose func(YieldDiffer) error

func (c Compose) runPlan(task Task, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	task.compose(c, leftPath, left, rightPath, right, differ)
}

// Delegate returns a Plan that replaces the current Differ with the provided Differ for the remainder of the comparison.
func Delegate(differ Differ2) Plan

type delegateImpl struct {
	differ Differ2
}

func (d delegateImpl) runPlan(task Task, leftPath Path, left Entry, rightPath Path, right Entry, differ Differ2) {
	task.delegate(leftPath, left, rightPath, right, differ)
}