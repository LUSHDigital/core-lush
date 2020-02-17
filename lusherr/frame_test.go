package lusherr_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/LUSHDigital/core-lush/lusherr"
	"github.com/LUSHDigital/core/test"
)

var dir string

func TestMain(m *testing.M) {
	var err error
	dir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestLocate(t *testing.T) {
	internal := lusherr.NewInternalError(fmt.Errorf("inner"))
	wrapped := fmt.Errorf("wrapping: %w", lusherr.NewInternalError(fmt.Errorf("wrapped")))
	untyped := fmt.Errorf("hello world")
	type Test struct {
		name     string
		err      error
		expect   bool
		expected runtime.Frame
	}
	cases := []Test{
		{
			name:   "with re-originated inline error",
			err:    lusherr.Originate(fmt.Errorf("inline error")),
			expect: true,
			expected: runtime.Frame{
				Line:     38,
				File:     filepath.Join(dir, "frame_test.go"),
				Function: "github.com/LUSHDigital/core-lush/lusherr_test.TestLocate",
			},
		},
		{
			name:   "with re-originated typed error",
			err:    lusherr.Originate(internal),
			expect: true,
			expected: runtime.Frame{
				Line:     48,
				File:     filepath.Join(dir, "frame_test.go"),
				Function: "github.com/LUSHDigital/core-lush/lusherr_test.TestLocate",
			},
		},
		{
			name:   "with error wrapped origin",
			err:    wrapped,
			expect: true,
			expected: runtime.Frame{
				Line:     27,
				File:     filepath.Join(dir, "frame_test.go"),
				Function: "github.com/LUSHDigital/core-lush/lusherr_test.TestLocate",
			},
		},
		{
			name:   "with untyped error",
			err:    untyped,
			expect: false,
			expected: runtime.Frame{
				Line:     28,
				File:     filepath.Join(dir, "frame_test.go"),
				Function: "github.com/LUSHDigital/core-lush/lusherr_test.TestLocate",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			frame, ok := lusherr.Locate(c.err)
			if c.expect && !ok {
				t.Fatal("frame not found")
			}
			if !c.expect && ok {
				t.Fatal("frame found when none was expected")
			}
			if ok {
				test.Equals(t, c.expected.File, frame.File)
				test.Equals(t, c.expected.Function, frame.Function)
				test.Equals(t, c.expected.Line, frame.Line)
			}
		})
	}
}
