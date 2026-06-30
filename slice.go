package diff

import (
	"fmt"
	"iter"
	"math"
	"reflect"
)

// SliceKeyFunc is a function that takes an index and an element of a slice and returns a key for
// aligning elements between two slices. The returned key must comparable and may be nil.
// This is used to reduce the amount of work required by the Slice differ to determine
// which elements in the left and right slices should be compared to each other, especially in cases
// where slice order does not matter.
type SliceKeyFunc func(index int, elem any) (any, error)

type Slice struct {
	// ElemDiffer is the Differ used to compare the elements of the slice.
	// If nil, the Default Differ is used.
	ElemDiffer Differ

	// AllowLeftOnly indicates whether elements that are only present in the left slice should be ignored.
	AllowLeftOnly bool

	// AllowRightOnly indicates whether elements that are only present in the right slice should be ignored.
	AllowRightOnly bool

	// TreatNilAsEmpty indicates whether nil slices should be treated as empty slices.
	TreatNilAsEmpty bool

	// KeyFunc is a function used to determine the key for aligning elements.
	// If nil, MatchByIndex is used, which aligns elements by their index in the slice.
	//
	// The returned value must be comparable and must never be nil.
	KeyFunc SliceKeyFunc
}

// SliceKeyByIndex is the default SliceKeyFunc.  It aligns elements by their index in the slice.
func SliceKeyByIndex(index int, _ any) (any, error) {
	return index, nil
}

// SliceKeyByElement is a key function that aligns elements by their elements.
// This only works for slices of comparable types. If the element type is not comparable, an error is returned.
func SliceKeyByElement(_ int, elem any) (any, error) {
	if elem == nil {
		return nil, nil
	}
	elemType := reflect.TypeOf(elem)
	if !elemType.Comparable() {
		return nil, fmt.Errorf("element type %s is not comparable", elemType)
	}
	return elem, nil
}

// SliceKeyInefficient is a key function that aligns all entries of a slice to the same key.
// Use this only in cases where the slice is small, the order of elements does not matter, and
// where it is difficult to determine a key for aligning elements.
// This will result in a factorial number of comparisons being made.
func SliceKeyInefficient(_ int, _ any) (any, error) {
	return nil, nil
}

func (s Slice) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, fmt.Errorf("state must not be nil")
	}

	if left == nil || right == nil {
		return false, WrapError(state.Path, fmt.Errorf("left and right must not be nil"))
	}

	leftVal := reflect.ValueOf(left)
	rightVal := reflect.ValueOf(right)

	if leftVal.Kind() != reflect.Slice || rightVal.Kind() != reflect.Slice {
		return false, WrapError(state.Path, fmt.Errorf("left and right must be slices"))
	}

	if leftVal.IsNil() && rightVal.IsNil() {
		return true, nil
	}

	if !s.TreatNilAsEmpty && (leftVal.IsNil() || rightVal.IsNil()) {
		state.Different()
		return false, nil
	}

	elemDiffer := s.ElemDiffer
	if elemDiffer == nil {
		elemDiffer = Default{}
	}

	matchFunc := s.KeyFunc
	if matchFunc == nil {
		matchFunc = SliceKeyByIndex
	}

	type matchup struct {
		left  []int
		right []int
	}

	matches := make(map[any]*matchup)
	for i := 0; i < leftVal.Len(); i++ {
		elem := leftVal.Index(i).Interface()
		key, err := matchFunc(i, elem)
		if err != nil {
			return false, WrapError(state.Path, fmt.Errorf("error in MatchFunc at left index %d: %w", i, err))
		}
		if found, exists := matches[key]; !exists {
			matches[key] = &matchup{left: []int{i}}
		} else {
			found.left = append(found.left, i)
		}
	}

	for i := 0; i < rightVal.Len(); i++ {
		elem := rightVal.Index(i).Interface()
		key, err := matchFunc(i, elem)
		if err != nil {
			return false, WrapError(state.Path, fmt.Errorf("error in MatchFunc at right index %d: %w", i, err))
		}
		if found, exists := matches[key]; !exists {
			matches[key] = &matchup{right: []int{i}}
		} else {
			found.right = append(found.right, i)
		}
	}

	allSame := true
	reportRightOnly := func(i int) {
		if !s.AllowRightOnly {
			rightElem := rightVal.Index(i).Interface()
			state.RightOnly(NewSliceUnmatchedKey(i), rightElem)
			allSame = false
		}
	}
	reportLeftOnly := func(i int) {
		if !s.AllowLeftOnly {
			leftElem := leftVal.Index(i).Interface()
			state.LeftOnly(NewSliceUnmatchedKey(i), leftElem)
			allSame = false
		}
	}
	for _, match := range matches {
		switch {
		case len(match.left) == 0:
			for _, rightIndex := range match.right {
				reportRightOnly(rightIndex)
			}
			continue
		case len(match.right) == 0:
			for _, leftIndex := range match.left {
				reportLeftOnly(leftIndex)
			}
			continue
		case len(match.left) == 1 && len(match.right) == 1:
			leftIndex := match.left[0]
			rightIndex := match.right[0]
			leftElem := leftVal.Index(leftIndex).Interface()
			rightElem := rightVal.Index(rightIndex).Interface()
			same, err := state.DiffChild(NewSliceKey(leftIndex, rightIndex), leftElem, rightElem, elemDiffer)
			if err != nil {
				return false, WrapError(state.Path, fmt.Errorf("error comparing left index %d and right index %d: %w", leftIndex, rightIndex, err))
			}
			if !same {
				allSame = false
			}
			continue
		}

		// Compare matched elements in a way that minimizes the number of diffs reported.
		var largerName, smallerName string
		var larger, smaller []int
		var largerVal, smallerVal reflect.Value
		var allowLargerOnly, allowSmallerOnly bool
		var reportLargerOnly func(int)
		var diffChild func(int, int) error
		if len(match.left) >= len(match.right) {
			largerName = "left"
			smallerName = "right"
			larger = match.left
			smaller = match.right
			largerVal = leftVal
			smallerVal = rightVal
			allowLargerOnly = s.AllowLeftOnly
			allowSmallerOnly = s.AllowRightOnly
			reportLargerOnly = reportLeftOnly
			diffChild = func(largerIndex, smallerIndex int) error {
				leftElem := leftVal.Index(largerIndex).Interface()
				rightElem := rightVal.Index(smallerIndex).Interface()
				same, err := state.DiffChild(NewSliceKey(largerIndex, smallerIndex), leftElem, rightElem, elemDiffer)
				if err != nil {
					allSame = false
					return err
				}
				if !same {
					allSame = false
				}
				return nil
			}
		} else {
			largerName = "right"
			smallerName = "left"
			larger = match.right
			smaller = match.left
			largerVal = rightVal
			smallerVal = leftVal
			allowLargerOnly = s.AllowRightOnly
			allowSmallerOnly = s.AllowLeftOnly
			reportLargerOnly = reportRightOnly
			diffChild = func(largerIndex, smallerIndex int) error {
				leftElem := leftVal.Index(smallerIndex).Interface()
				rightElem := rightVal.Index(largerIndex).Interface()
				same, err := state.DiffChild(NewSliceKey(smallerIndex, largerIndex), leftElem, rightElem, elemDiffer)
				if err != nil {
					allSame = false
					return err
				}
				if !same {
					allSame = false
				}
				return nil
			}
		}
		type sizeMatchup struct {
			larger, smaller []int
		}
		newSizeMatchup := func(larger, smaller []int) sizeMatchup {
			return sizeMatchup{
				larger:  append([]int(nil), larger...),
				smaller: append([]int(nil), smaller...),
			}
		}
		permuteSlice := func(in []int) iter.Seq[[]int] {
			return func(yield func([]int) bool) {
				if len(in) == 0 {
					yield([]int{})
					return
				}

				working := append([]int(nil), in...)

				var build func(start int) bool
				build = func(start int) bool {
					if start == len(working)-1 {
						return yield(append([]int(nil), working...))
					}

					for i := start; i < len(working); i++ {
						working[start], working[i] = working[i], working[start]
						stop := !build(start + 1)
						working[start], working[i] = working[i], working[start]
						if stop {
							return false
						}
					}
					return true
				}

				_ = build(0)
			}
		}
		yieldSizeMatchups := func(larger, smaller []int) iter.Seq[sizeMatchup] {
			return func(yield func(sizeMatchup) bool) {
				for largerPerm := range permuteSlice(larger) {
					if !yield(newSizeMatchup(largerPerm, smaller)) {
						return
					}
				}
			}
		}
		indirect := func(inVal reflect.Value, inIndexes []int) []any {
			out := make([]any, len(inIndexes))
			for i, index := range inIndexes {
				out[i] = inVal.Index(index).Interface()
			}
			return out
		}
		var bestMatchup sizeMatchup
		bestDiffCount := math.MaxInt
		for matchup := range yieldSizeMatchups(larger, smaller) {
			largerSubset := indirect(largerVal, matchup.larger)
			smallerSubset := indirect(smallerVal, matchup.smaller)
			var counter Counter
			_, err := Diff(largerSubset, smallerSubset, Slice{
				ElemDiffer:     elemDiffer,
				AllowLeftOnly:  allowLargerOnly,
				AllowRightOnly: allowSmallerOnly,
				KeyFunc:        SliceKeyByIndex,
			}, &counter)
			if err != nil {
				return false, WrapError(state.Path, fmt.Errorf("error comparing %s and %s slices: %w", largerName, smallerName, err))
			}
			if counter.Total() < bestDiffCount {
				bestDiffCount = counter.Total()
				bestMatchup = matchup
			}
			if bestDiffCount == 0 {
				break
			}
		}
		for i, largerIndex := range bestMatchup.larger {
			if i >= len(bestMatchup.smaller) {
				reportLargerOnly(largerIndex)
			} else {
				smallerIndex := bestMatchup.smaller[i]
				if err := diffChild(largerIndex, smallerIndex); err != nil {
					return false, WrapError(state.Path, fmt.Errorf("error comparing %s index %d and %s index %d: %w", largerName, largerIndex, smallerName, smallerIndex, err))
				}
			}
		}
	}

	return allSame, nil
}
