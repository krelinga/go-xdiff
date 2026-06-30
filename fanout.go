package diff

type Fanout []Reporter

func (f Fanout) Push(key Key, left, right any) {
	for _, reporter := range f {
		reporter.Push(key, left, right)
	}
}

func (f Fanout) Pop() {
	for _, reporter := range f {
		reporter.Pop()
	}
}

func (f Fanout) LeftOnly(key Key, left any) {
	for _, reporter := range f {
		reporter.LeftOnly(key, left)
	}
}

func (f Fanout) RightOnly(key Key, right any) {
	for _, reporter := range f {
		reporter.RightOnly(key, right)
	}
}

func (f Fanout) Different() {
	for _, reporter := range f {
		reporter.Different()
	}
}
