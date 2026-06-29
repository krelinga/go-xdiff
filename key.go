package diff

import fmt

// FieldKey represents a field in a struct.
type FieldKey struct {
	// The name of the struct field.  This should always be an upper-case (i.e. exported) name.
	Name string
}

func (fk FieldKey) DiffKey() string {
	return fmt.Sprintf("struct field: %s", fk.Name)
}

func NewFieldKey(name string) FieldKey {
	return FieldKey{
		Name: name,
	}
}

// MapKey represents a key in a map.
type MapKey struct {
	// Must be comparable.
	Key any
}

func (fk FieldKey) DiffKey() string {
	return fmt.Sprintf("map key: %v", fk.Key)
}

func NewMapKey(key any) MapKey {
	return MapKey{
		Key: key,
	}
}

// SliceUnmatchedKey represents a slice index that only existed in one side of the diff.
type SliceUnmatchedKey struct {
	Index int
}

func NewSliceUnmatchedKey(index int) SliceUnmatchedKey {
	return SliceUnmatchedKey{
		Index: index,
	}
}

func (suk SliceUnmatchedKey) DiffKey() string {
	return fmt.Sprintf("slice unmatched index: %d", suk.Index)
}

// SliceKey represents a pair of matched elements in left and right slices, although not necessairly at the same index.
type SliceKey {
	LeftIndex, RightIndex int
}

func NewSliceKey(leftIndex, rightIndex int) SliceKey {
	return SliceKey{
		LeftIndex: leftIndex,
		RightIndex: rightIndex,
	}
}

// RootKey is used as a Key for the root of the diff operation.
type RootKey struct {}

func (_ RootKey) DiffKey() string {
	return "root"
}
