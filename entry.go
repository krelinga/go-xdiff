package diff

type Entry struct {
	ok bool
	a any
}

func (e Entry) Ok() bool {
	return e.ok
}

func (e Entry) Get() (any, bool) {
	return e.a, e.ok
}

func (e Entry) Must() any {
	if !e.ok {
		panic("entry invalid")
	}
	return e.a
}

func NewEntry(a any) Entry {
	return Entry{ok: true, a: a}
}

func NoEntry() Entry {
	return Entry{}
}