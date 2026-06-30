package diff

// Counter is a Reporter that counts the number of differences, left-only elements, and right-only elements.
type Counter struct {
	NumDifferent int
	NumLeftOnly  int
	NumRightOnly int
}

func (c *Counter) Different() {
	c.NumDifferent++
}

func (c *Counter) LeftOnly(_ Key, _ any) {
	c.NumLeftOnly++
}

func (c *Counter) RightOnly(_ Key, _ any) {
	c.NumRightOnly++
}

func (_ *Counter) Push(_ Key, _ string, _, _ any) {}
func (_ *Counter) Pop()                 {}

func (c *Counter) Total() int {
	return c.NumDifferent + c.NumLeftOnly + c.NumRightOnly
}
