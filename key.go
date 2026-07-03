package diff

import "fmt"

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

func (mk MapKey) DiffKey() string {
	return fmt.Sprintf("map key: %v", mk.Key)
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

// SliceKey represents a pair of matched elements in left and right slices, although not necessarily at the same index.
type SliceKey struct {
	LeftIndex, RightIndex int
}

func NewSliceKey(leftIndex, rightIndex int) SliceKey {
	return SliceKey{
		LeftIndex: leftIndex,
		RightIndex: rightIndex,
	}
}

func (sk SliceKey) DiffKey() string {
	return fmt.Sprintf("slice pair: left index %d, right index %d", sk.LeftIndex, sk.RightIndex)
}

// RootKey is used as a Key for the root of the diff operation.
type RootKey struct {}

func (_ RootKey) DiffKey() string {
	return "root"
}

// PointerKey represents a pointer dereference in the diff path.
type PointerKey struct {}

func (_ PointerKey) DiffKey() string {
	return "pointer dereference"
}

type TransformKey struct {
	// The name of the transform.
	Name string
}

func (tk TransformKey) DiffKey() string {
	return fmt.Sprintf("transform: %s", tk.Name)
}

func NewTransformKey(name string) TransformKey {
	return TransformKey{
		Name: name,
	}
}

type composeKey struct {}

func (_ composeKey) DiffKey() string {
	return "compose"
}