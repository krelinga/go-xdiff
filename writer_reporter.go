package diff

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type writerLevel struct {
	name string
	key Key
	left, right any
}

type WriterReporter struct {
	// Should be set at construction time.
	Writer io.Writer

	lastOutputPath Path
	levels []writerLevel
}

func (r *WriterReporter) currentPath() Path {
	path := make(Path, len(r.levels))
	for i, level := range r.levels {
		path[i] = level.key
	}
	return path
}

func (r *WriterReporter) writeIndent() {
	const indent = "  "
	for i := 0; i < len(r.levels); i++ {
		currentLevel := r.levels[i]
		if i > 0 {
			io.WriteString(r.Writer, indent)
		}
		currentKey := currentLevel.key
		var lastKey Key
		if i < len(r.lastOutputPath) {
			lastKey = r.lastOutputPath[i]
		}
		if currentKey == lastKey {
			continue
		}

		extraIndent := strings.Repeat(indent, i)
		var typePart string
		leftType := reflect.TypeOf(currentLevel.left)
		rightType := reflect.TypeOf(currentLevel.right)
		if leftType == rightType {
			typePart = leftType.String()
		} else {
			typePart = fmt.Sprintf("%s vs %s", leftType, rightType)
		}
		fmt.Fprintf(r.Writer, "%s %s %s \n%s", currentKey.DiffKey(), currentLevel.name, typePart, extraIndent)
	}
	io.WriteString(r.Writer, " ")
	r.lastOutputPath = r.currentPath()
}

func (r *WriterReporter) checkWriter() {
	if r.Writer == nil {
		panic("WriterReporter.Writer must not be nil")
	}
}

func (r *WriterReporter) Push(key Key, name string, left, right any) {
	r.checkWriter()
	r.levels = append(r.levels, writerLevel{
		name:      name,
		key:       key,
		left: left,
		right: right,
	})
}

func (r *WriterReporter) Pop() {
	r.checkWriter()
	if len(r.levels) == 0 {
		panic("WriterReporter.Pop called without a matching Push")
	}
	r.levels = r.levels[:len(r.levels)-1]
}

func (r *WriterReporter) LeftOnly(key Key, left any) {
	r.checkWriter()
	r.writeIndent()
	fmt.Fprintf(r.Writer, "- %s %v\n", key.DiffKey(), left)
}

func (r *WriterReporter) RightOnly(key Key, right any) {
	r.checkWriter()
	r.writeIndent()
	fmt.Fprintf(r.Writer, "+ %s %v\n", key.DiffKey(), right)
}

func (r *WriterReporter) Different() {
	r.checkWriter()
	top := r.levels[len(r.levels)-1]
	leftStr := fmt.Sprintf("%v", top.left)
	rightStr := fmt.Sprintf("%v", top.right)
	divider := strings.Repeat("-", max(len(leftStr), len(rightStr)))
	r.writeIndent()
	fmt.Fprintf(r.Writer, "! %v\n", top.left)
	r.writeIndent()
	fmt.Fprintf(r.Writer, "! %s\n", divider)
	r.writeIndent()
	fmt.Fprintf(r.Writer, "! %v\n", top.right)
}
