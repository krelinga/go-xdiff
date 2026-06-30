package diff_test

import (
	"os"

	diff "github.com/krelinga/go-xdiff"
)

func ExampleWriterReporter() {
	type FooStruct struct {
		Name       string
		Age        int
		Properties map[string]string
	}

	left := &FooStruct{
		Name: "Alice",
		Age:  30,
		Properties: map[string]string{
			"city": "New York",
			"job":  "Engineer",
		},
	}

	right := &FooStruct{
		Name: "Alice",
		Age:  31,
		Properties: map[string]string{
			"city":  "San Francisco",
			"job":   "Engineer",
			"hobby": "Photography",
		},
	}

	same, err := diff.Diff(left, right, diff.Default{}, &diff.WriterReporter{Writer: os.Stdout})
	if err != nil {
		panic(err)
	}

	if same {
		println("The two structures are the same.")
	} else {
		println("The two structures are different.")
	}

	// Output:
}
